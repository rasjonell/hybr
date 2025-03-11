package services

func GetRegisteredServices() []HybrService {
	mu.RLock()
	defer mu.RUnlock()

	services := make([]HybrService, 0, len(registry))
	for _, s := range registry {
		services = append(services, s)
	}
	return services
}

func findPort(varDef []*VariableDefinition) string {
	var port string
	for _, def := range varDef {
		if def.Name == "PORT" {
			port = def.Value
		}
	}

	return port
}
