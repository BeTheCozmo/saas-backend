package main

import uller "uller/src"

func main() {
  manager := uller.New()
  manager.CreateModules()
  manager.ConfigureModules()
  manager.Run()
}
