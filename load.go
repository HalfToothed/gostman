package main

func load(data Request, model *Model) {

	model.nameField.SetValue(data.Name)
	model.urlField.SetValue(data.URL)
	model.methodField.SetValue(data.Method)
	model.tabContent[0].SetValue(data.Body)
	model.tabContent[1].SetValue(data.QueryParams)
	model.tabContent[2].SetValue(data.Headers)
	model.response = data.Response
}
