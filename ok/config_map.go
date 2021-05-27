package ok

type Map map[string]interface{}

type ConfigMap Map

func (m ConfigMap) Debug() bool {
	debug, debugExists := m["debug"]
	return debugExists && debug.(bool)
}

func (m ConfigMap) AddSkipTool(toolName string) {
	m.ensureStringArray("skip")
	m["skip"] = append(m["skip"].([]string), toolName)
}

func (m ConfigMap) AddToolSort(toolName string) {
	m.ensureStringArray("tool_sort")
	m["tool_sort"] = append(m["tool_sort"].([]string), toolName)
}

func (m ConfigMap) SkipTools() (skipTools []string) {
	m.ensureStringArray("skip")
	return
}

func (m ConfigMap) ToolSort() []string {
	m.ensureStringArray("tool_sort")
	return m["tool_sort"].([]string)
}

func (m ConfigMap) ToolConfig(toolName string) (map[string]interface{}, bool) {
	toolConfig, toolExists := m[toolName]
	return toolConfig.(map[string]interface{}), toolExists
}

func (m ConfigMap) ensureStringArray(key string) {
	if _, exists := m[key]; !exists {
		m[key] = []string{}
	}
}
