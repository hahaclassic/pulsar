# pulsar

>I once set out to walk a hundred thousand steps in a single day. By the sixty-three-thousandth, my heart ached fiercely,
though I have devoted my life to sport and movement. For the first time, I felt true, unshakable fear. In such moments, one realizes: a human is not a machine, and a machine, in turn, is not human. Yet what if — just for a fleeting instant — we could endow a soulless machine with a fragment of humanity?

`Pulsar` is a cardiometer for your computer — a playful tool that lets you see the “heartbeat” of your CPU in real time.

```
         |                   |                   |               
        /|                  /|                  /|               
       / |   |             / |   |             / |   |           
 _ _  /  |  /| /\_ _ _ _  /  |  /| /\_ _ _ _  /  |  /| /\______  
‾ ‾ \/   | / |/   ‾ ‾ ‾ \/   | / |/   ‾ ‾ ‾ \/   | / |/       
         |/                  |/                  |/              
         |                   |                   |             
cpu: 60.2%      bpm: 129      
```

## Features
- Visualizes CPU usage as a heartbeat.
- Shows a real-time BPM (beats per minute) based on system load.
- Supports terminal resizing and keeps the waveform aligned.
- Graceful shutdown (via `Ctrl+C` or `SIGTERM`).
- Lightweight and fun way to monitor your computer.

## Platform Support

- Linux only (requires `/proc/stat` for CPU usage)

## Installation

### Clone the repository
```bash
git clone https://github.com/hahaclassic/pulsar.git
cd pulsar
```

### Build
```bash
make
```
or
```bash
go build -o pulsar ./cmd/pulsar/main.go
```

### Enjoy
``` 
./pulsar
```

## Optional: CPU Stress Test

To see the heartbeat react to high CPU usage, you can run a small stress test:

```bash
go run ./tools/cpu_stress.go -n 12
```
or 
```bash
make stress-test
```

This launches multiple goroutines (n is number of goroutines) to fully load your CPU, allowing you to observe the Pulsar reacting with higher BPM values.


## License
Licensed under [MIT License](./LICENSE).
