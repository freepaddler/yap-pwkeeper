package models

var Store Wallet = Wallet{
	Login: []LoginItem{
		{
			Id:       "1",
			Name:     "name1",
			Login:    "login1",
			Password: "password1",
			MetaQ: MetaQ{
				MetaItem{
					key:   "key1",
					value: "value1",
				},
				MetaItem{
					key:   "key2",
					value: "value2",
				},
			},
		},
		{
			Id:       "2",
			Name:     "name2",
			Login:    "login2",
			Password: "password2",
			MetaQ: MetaQ{
				MetaItem{
					key:   "key1",
					value: "value1",
				},
				MetaItem{
					key:   "key3",
					value: "value3",
				},
			},
		},
	},
	Note: []NoteItem{
		{
			Id:    "3",
			Name:  "name3",
			Text:  "this is a text3 text",
			MetaQ: nil,
		},
		{
			Id:   "4",
			Name: "name4",
			Text: "this is a text4 text",
			MetaQ: MetaQ{
				MetaItem{
					key:   "metakey",
					value: "metavalue",
				},
			},
		},
	},
}
