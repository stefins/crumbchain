package main

import (
  "testing"

  "github.com/larissareyes/crumbchain/crumb"
  "github.com/larissareyes/crumbchain/crumbjoiner"
)

func TestCrumbing(t *testing.T)  {
  crumb.Crumber("4_188139517687892532.mp4")
}

func TestJoining(t *testing.T) {
  crumbjoiner.Joiner("4_188139517687892532.mp4-crumb")
}
