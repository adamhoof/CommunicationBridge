package telegram

import tb "gopkg.in/telebot.v3"

const AllRoomsCommand = "/all_rooms"
const OfficeToysCommand = "/office_toys"
const BedroomToysCommand = "/bedroom_toys"
const TableLampOCommand = "/lamp_office"
const CeilingLightOCommand = "/ceil_light_office"
const TableLampBCommand = "/lamp_bedroom"
const ShadesBCommand = "/shades_bedroom"

func RoomCommands() []tb.Command {
	return []tb.Command{
		{Text: AllRoomsCommand, Description: "All rooms"},
		{Text: OfficeToysCommand, Description: "Office Toys"},
		{Text: BedroomToysCommand, Description: "Bedroom Toys"},
	}
}

func OfficeToysCommands() []tb.Command {
	return []tb.Command{
		{Text: TableLampOCommand, Description: ""},
		{Text: CeilingLightOCommand, Description: ""},
	}
}

func BedroomToysCommands() []tb.Command {
	return []tb.Command{
		{Text: TableLampBCommand, Description: ""},
		{Text: ShadesBCommand, Description: ""},
	}
}
