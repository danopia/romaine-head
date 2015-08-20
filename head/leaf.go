package head

import (
  "os/exec"
)

type Leaf struct {
  Running bool
  Secret  string
  Anchor  *exec.Cmd
}
