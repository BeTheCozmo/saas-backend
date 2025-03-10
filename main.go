package main

import manager "uller/src"

func main() {
  manager := manager.New()
  manager.CreateModules()
  manager.ConfigureModules()
  manager.Run()
}
