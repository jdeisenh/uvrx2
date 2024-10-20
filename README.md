Experimentelle Abfrage einer UVR Generation X2 Heizungsstuerung.

Disclaimer:
Diese Software ist kein Produkt der Firma Technische Alternative und hat keinerlei Freigabe oder Unterstützung dieser Firma. 
Diese Software dient Lernzwecken. Sie kann Schäden an ihrer Heizunganlage und in Folge darüber hinaus verursachen, für die niemand, insbesondere nicht ich, haften wird. Die Gefahr beim reinen Auslesen halten ich allerdings für gering. 


Basiert auf can/canopen library from 

Mit dieser Software können Mess- un Konfigurationswerte aus einer UVR X2 Heizungssteuerung der 
Firma Technische Alternative direkt über CAN ausgelesen werden.
Dafür ist entsprechende CAN-Hardware nötig, die sehr preisgünstig erhältlich ist. 
Ein CMI-Modul ist dafür nicht notwendig. 

Funktionalität:
Auslesen aller Mess- und Konfigurationswerte, die auch per CMI erreichbar sind.
Export nach prometheus, zur Überwachung und Visualisierung z.B. mittels grafana

Ein Wort vorneweg: Da es sich hier um Universal-Regler handelt, hängt die Interpretation 
der ausgelesenen Werte von der Programmierung ab. Die Werte aus der Programmierung muss von Hand übertragen werden. Wer es einfach und zuverlässig haben will, kauft weiterhin am besten das CMI-Modul

Hardware: 
Ich verwende
- Rasperry Zero W2
- waveshare CAN-Hat
- DC/DC Wandler zur Stromversorgung des Raspi direkt aus der UVR
- Gehäuse, in die alle drei Komponenten eingebaut werden können

Unterstützte Geräte:
Getestet ausschliesslich mit UVR 1610. Die Geräte aus derselben Generation ... sollten jedoch sehr ähnlich funktioieren. 

TODO:
Schreibzugriff

Das Format der Schreibzugriffe an die Steuerung konnte ich bisher nicht reverse-engineeren. Dazu wäre ein Belauschen des CAN-Busses während Schreibzugriffe notwendig.
Ich wäre froh, wenn dies jemand für mich erledigen könnte, oder mir für eine begrenzte Zeit ein CMI-Modul zur Verfüngung stellt. 

MQTT Unterstützung zum Setzen des Betriebsmodus

