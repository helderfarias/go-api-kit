package discovery

type ServiceDiscoveryRegister interface {
	Register()

	UnRegister()
}
