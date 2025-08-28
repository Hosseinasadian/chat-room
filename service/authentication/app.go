package authentication

type Application struct {
}

func Setup() Application {
	return Application{}
}

func (app Application) Start() {

}

// development
// config.yaml,dockerfile,docker-compose,...

// cmd
// command line to start(serve) APPS
