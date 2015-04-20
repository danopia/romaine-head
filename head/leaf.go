package head

import (
  "os/exec"
)

type Leaf struct {
  Running bool
  Port    int
  Anchor  *exec.Cmd
}
