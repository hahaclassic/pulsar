# pulsar

>I once set out to walk a hundred thousand steps in a single day. By the sixty-three-thousandth, my heart ached fiercely,
though I have devoted my life to sport and movement. For the first time, I felt true, unshakable fear. In such moments, one realizes: a human is not a machine, and a machine, in turn, is not human. Yet what if — just for a fleeting instant — we could endow a soulless machine with a fragment of humanity?

`Pulsar` is a cardiometer for your computer — a playful tool that lets you see the “heartbeat” of your CPU in real time.

```
         |                  |                  |               
        /|                 /|                 /|               
       / |   |            / |   |            / |   |           
____  /  |  /| /\______  /  |  /| /\______  /  |  /| /\______  
    \/   | / |/        \/   | / |/        \/   | / |/       
         |/                 |/                 |/              
         |                  |                  |             
cpu: 60.2%      bpm: 129      
```

## Features
- Visualizes CPU usage as a heartbeat.
- Shows a real-time BPM (beats per minute) based on system load.
- Lightweight and fun way to monitor your computer.

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

## License
Licensed under [MIT License](./LICENSE).
