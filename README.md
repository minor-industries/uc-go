# uc-go

A collection of tinygo/microcontroller sketches and prototypes.

## Bikelights

A small, low-cost controller for ws2812 LED strips.

![bikeleds_board.jpg](https://minor-industries.sfo2.digitaloceanspaces.com/hw/bikeleds_board.jpg "bikeleds_board.jpg")

![bikeleds_with_case.jpg](https://minor-industries.sfo2.digitaloceanspaces.com/hw/bikeleds_with_case.jpg "bikeleds_with_case.jpg")

Features:

- SAMD51 microcontroller: 120 MHz, hardware floating point.
- SPI-based driver: This uses the SAMD SPI hardware peripheral and interrupts to clock out the ws2812 signal, allowing
  for animation code and other interrupts to run concurrently.
- [IR remote](https://adafruit.com/product/389) control: change animations, brightness, etc.
- Flash filesystem: persist animation settings, brightness, etc. between restarts.
- 3D printed case with JST SM connector: JST SM connectors have no wire to board connectors, so we instead use our own 3d printed connector which keeps the size small and is easy to plug and unplug.
- USB Bootloader: custom bootloaders can be uploaded with J-Link.
- Easy to assemble.

Hardware is available [here](https://github.com/minor-industries/hardware/tree/main/bikeleds). 

