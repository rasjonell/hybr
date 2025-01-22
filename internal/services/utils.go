package services

func GetRegisteredServices() []Service {
	mu.RLock()
	defer mu.RUnlock()

	services := make([]Service, 0, len(registry))
	for _, s := range registry {
		services = append(services, s)
	}
	return services
}
