package service

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/docker/docker/pkg/term"
	"github.com/klovercloud-ci-cd/terminal/config"
	"gopkg.in/igm/sockjs-go.v2/sockjs"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/kubernetes/pkg/util/interrupt"
	"log"
	"net/http"
)

func (self TerminalSockjs) Read(p []byte) (int, error) {
	var reply string
	var msg map[string]uint16
	reply, err := self.conn.Recv()
	if err != nil {
		return 0, err
	}
	if err := json.Unmarshal([]byte(reply), &msg); err != nil {
		return copy(p, reply), nil
	} else {
		self.sizeChan <- &remotecommand.TerminalSize{
			msg["cols"],
			msg["rows"],
		}
		return 0, nil
	}
}

func (self TerminalSockjs) Write(p []byte) (int, error) {
	err := self.conn.Send(string(p))
	return len(p), err
}

type TerminalSockjs struct {
	conn      sockjs.Session
	sizeChan  chan *remotecommand.TerminalSize
	context   string
	companyId string
	namespace string
	pod       string
	container string
}

// queue
func (self *TerminalSockjs) Next() *remotecommand.TerminalSize {
	size := <-self.sizeChan
	beego.Debug(fmt.Sprintf("terminal size to width: %d height: %d", size.Width, size.Height))
	return size
}

// Processing input and output interact with sockjs
func Handler(t *TerminalSockjs, cmd string) error {
	config := config.KubeConfig
	groupversion := schema.GroupVersion{
		Group:   "",
		Version: "v1",
	}
	config.GroupVersion = &groupversion
	config.APIPath = "/api"
	config.ContentType = runtime.ContentTypeJSON
	config.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: scheme.Codecs}
	restclient, err := rest.RESTClientFor(config)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	fn := func() error {
		req := restclient.Post().
			Resource("pods").
			Name(t.pod).
			Namespace(t.namespace).
			SubResource("exec").
			Param("container", t.container).
			Param("stdin", "true").
			Param("stdout", "true").
			Param("stderr", "true").
			Param("command", cmd).Param("tty", "true")
		req.VersionedParams(
			&v1.PodExecOptions{
				Container: t.container,
				Command:   []string{},
				Stdin:     true,
				Stdout:    true,
				Stderr:    true,
				TTY:       true,
			},
			scheme.ParameterCodec,
		)
		executor, err := remotecommand.NewSPDYExecutor(
			config, http.MethodPost, req.URL(),
		)
		if err != nil {
			return err
		}
		return executor.Stream(remotecommand.StreamOptions{
			Stdin:             t,
			Stdout:            t,
			Stderr:            t,
			Tty:               true,
			TerminalSizeQueue: t,
		})
	}
	inFd, _ := term.GetFdInfo(t.conn)
	state, err := term.SaveState(inFd)
	return interrupt.Chain(nil, func() {
		term.RestoreTerminal(inFd, state)
	}).Run(fn)
}

// Implement http.handler interface to get input parameters
func (self TerminalSockjs) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	context := r.FormValue("context")
	//token := r.FormValue("token")
	pod := r.FormValue("pod")
	namespace := r.FormValue("namespace")
	container := r.FormValue("container")
	companyId := r.FormValue("companyId")
	token := r.FormValue("token")

	if config.EnableAuthentication {
		userResourcePermission, err := GetUserResourcePermissionFromBearerToken(token)
		if err != nil {
			log.Println(err.Error())
			w.Write([]byte("[ERROR]:" + err.Error()))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if err := checkAuthority(userResourcePermission, string(APPLICATION), "", string(READ)); err != nil {
			if err != nil {
				w.Write([]byte("[ERROR]:" + err.Error()))
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		}
		companyId = userResourcePermission.Metadata.CompanyId
	}

	Sockjshandler := func(session sockjs.Session) {
		t := &TerminalSockjs{session, make(chan *remotecommand.TerminalSize),
			context, companyId, namespace, pod, container}
		if err := Handler(t, "/bin/bash"); err != nil {
			log.Println(err.Error())
			beego.Error(err)
			beego.Error(Handler(t, "/bin/sh"))
		}
	}

	sockjs.NewHandler("/terminal/ws", sockjs.DefaultOptions, Sockjshandler).ServeHTTP(w, r)
}
