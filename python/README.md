# 1 Billion Row Challenge (1BRC) - Python

Attempting to do the challenge in python was kinda rough. With a starting time of 14 minutes I opted to start writing out multiple versions and start running multiple versions. On my Mac this did end up heating up the Mac which is probably what caused the slower times in the later versions. As such the times listed will only be for the times captured in PC.

With the attempts I also opted not to use pypy as the JIT compiler as I'm mostly searching for what could cause slows downs or speed ups in pythons. This does mean that all the times could be potentially lower.

Something else that made the challenge harder to tackle was trying to run the profiler and visualize what exactly is taking all the time for execution. Setting up py-spy didn't go well on Mac and cProfile doesn't provide much information. Will be re-running on PC to check if more information is provided or if I have better luck using py-spy or pyinstrument (TODO: Update this block)

## Setup, Run and Profile

```bash
python3 -m venv .venv
source .venv/bin/activate
pip3 install -r requirements.txt
python3 main.py
snakeviz cpu.prof
```