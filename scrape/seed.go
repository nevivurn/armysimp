package scrape

type seedGeneration struct {
	title    string
	channels []seedChannel
}

type seedChannel struct {
	title     string
	youtubeID string
	favorite  bool
}

var (
	seedData = []seedGeneration{
		{
			"Official", []seedChannel{
				{"Hololive JP", "UCJFZiqLMntJufDCHc6bQixg", false},
			},
		},
		{
			"Gen 0", []seedChannel{
				{"Tokino Sora", "UCp6993wxpyDPHUpavwDFqgg", false},
				{"AZKi", "UC0TXe_LYZ4scaW2XMyi5_kw", false},
				{"Roboco san", "UCDqI2jOz0weumE8s7paEk6g", false},
				{"Sakura Miko", "UC-hM6YJuNYVAmUWxeIr9FeA", false},
				{"Hoshimachi Suisei", "UC5CwaMl1eIgY8h02uZw7u8A", false},
			},
		},
		{
			"Gen 1", []seedChannel{
				{"Shirakami Fubuki", "UCdn5BQ06XqgXoAxIhbqw5Rg", false},
				{"Natsuiro Matsuri", "UCQ0UDLQCjY0rmuxCDE38FGg", false},
				{"Yozora Mel", "UCD8HOxPs4Xvsm8H0ZxXGiBw", false},
				{"Akai Haato", "UC1CfXB_kRs3C-zaeTG3oGyg", false},
				{"Aki Rosenthal", "UCFTLzh12_nrtzqBPsTCqenA", false},
			},
		},
		{
			"Gen 2", []seedChannel{
				{"Minato Aqua", "UC1opHUrw8rvnsadT-iGp7Cg", false},
				{"Yuzuki Choco", "UC1suqwovbL1kzsoaZgFZLKg", false},
				{"Nakiri Ayame", "UC7fk0CB07ly8oSl0aqKkqFg", false},
				{"Murasaki Shion", "UCXTpFs_3PqI41qX2d9tL2Rw", false},
				{"Oozora Subaru", "UCvzGlP9oQwU--Y0r9id_jnA", true},
			},
		},
		{
			"Gamers", []seedChannel{
				{"Ookami Mio", "UCp-5t9SrOQwXMU7iIjQfARg", true},
				{"Nekomata Okayu", "UCvaTdHTWBGv3MKj3KVqJVCw", true},
				{"Inugami Korone", "UChAnqc_AY5_I3Px5dig3X1Q", false},
			},
		},
		{
			"Gen 3", []seedChannel{
				{"Shiranui Flare", "UCvInZx9h3jC2JzsIzoOebWg", true},
				{"Shirogane Noel", "UCdyqAaZDKHXg4Ahi7VENThQ", false},
				{"Houshou Marine", "UCCzUftO8KOVkV4wQG1vkUvg", true},
				{"Usada Pekora", "UC1DCedRgGHBdm81E1llLhOQ", false},
				{"Uruha Rushia", "UCl_gCybOJRIgOXw6Qb4qJzQ", false},
			},
		},
		{
			"Gen 4", []seedChannel{
				{"Amane Kanata", "UCZlDXzGoo7d44bwdNObFacg", false},
				{"Kiryu Coco", "UCS9uQI-jC3DE0L4IpXyvr6w", true},
				{"Tsunomaki Watame", "UCqm3BQLlJfvkTsX_hvm0UmA", false},
				{"Tokoyami Towa", "UC1uv2Oq6kNxgATlCiez59hw", false},
				{"Himemori Luna", "UCa9Y57gfeY0Zro_noHRVrnw", false},
			},
		},
		{
			"ID Gen 1", []seedChannel{
				{"Ayunda Risu", "UCOyYb1c43VlX9rc_lT6NKQw", false},
				{"Moona Hoshinova", "UCP0BspO_AMEe3aQqqpo89Dg", false},
				{"Airani Iofifteen", "UCAoy6rzhSf4ydcYjJw3WoVg", false},
			},
		},
		{
			"Gen 5", []seedChannel{
				{"Yukihana Lamy", "UCFKOVgVbGmX65RxO3EtH3iw", true},
				{"Momosuzu Nene", "UCAWSyEs_Io8MtpY3m-zqILA", false},
				{"Shishiro Botan", "UCUKD-uaobj9jiqB-VXt71mA", false},
				{"Omaru Polka", "UCK9V2B22uJYu3N7eR_BT9QA", false},
				// TT
			},
		},
		{
			"EN Gen 1", []seedChannel{
				{"Mori Calliope", "UCL_qhgtOy0dy1Agp8vkySQg", true},
				{"Takanashi Kiara", "UCHsx4Hqa-1ORjQTh9TYDhww", false},
				{"Ninomae Ina'nis", "UCMwGHR0BTZuLsmjY_NT5Pwg", true},
				{"Gawr Gura", "UCoSrY_IQQVpmIRZ9Xf-y93g", false},
				{"Watson Amelia", "UCyl1z3jo3XHR1riLFKG5UAg", false},
			},
		},
		{
			"ID Gen 2", []seedChannel{
				{"Kureiji Ollie", "UCYz_5n-uDuChHtLo7My1HnQ", false},
				{"Anya Melfissa", "UC727SQYUvx5pDDGQpTICNWg", false},
				{"Pavolia Reine", "UChgTyjG-pdNvxxhdsXfHQ5Q", true},
			},
		},
	}
)
