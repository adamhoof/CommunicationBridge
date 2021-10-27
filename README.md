# RPICommunicationBridge
#### The purpose of this device is connecting multiple services together such as TelegramBot, MQTT broker and PostgreSQL handlers. 

- The brain -> RPI 4.
- Program functions -> Written in Go. Receive commands from user over TelegramBot API, RPI sends message to MQTT broker (running on remote SSH server), on related topic to an end device controller, such as lamp controller, to do some task.
This method reduces weaker device controller (i.e. ESP32) load, which only has to do given task and report its status, leaving most of the processing to RPI, which then writes this status to PostgreDB. 

