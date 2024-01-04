package main

import (
	"bytes"
	"testing"
)

type myWriterObj struct {
	data bytes.Buffer
}

func (m *myWriterObj) Write(p []byte) (n int, err error) {
	return m.data.Write(p)
}

func (m *myWriterObj) GetString() string {
	return m.data.String()
}

func (m *myWriterObj) Clean() {
	m.data = bytes.Buffer{}
}

func Test_app(t *testing.T) {
	myWriterInstance := myWriterObj{}
	tests := []struct {
		name string
		conf *Configuration
		want string
	}{
		{name: "searchWord_1",
			conf: &Configuration{
				afterFlag:      0,
				beforeFlag:     0,
				contextFlag:    0,
				countFlag:      false,
				ignoreCaseFlag: true,
				invertFlag:     false,
				fixedFlag:      false,
				lineNumFlag:    false,
				regexPattern:   "финики",
				files:          []string{"./test_files/grocery_list.txt"},
			},
			want: "финики 7 кг.\n",
		},

		{name: "afterFlag_1",
			conf: &Configuration{
				afterFlag:      3,
				beforeFlag:     0,
				contextFlag:    0,
				countFlag:      false,
				ignoreCaseFlag: true,
				invertFlag:     false,
				fixedFlag:      false,
				lineNumFlag:    false,
				regexPattern:   "финики",
				files:          []string{"./test_files/grocery_list.txt"},
			},
			want: "финики 7 кг.\n" +
				"конфеты 10 кг.\n" +
				"изотоп калия K40 1 г.\n" +
				"место на жестком диске +2 ТБ\n",
		},

		{name: "afterFlag_2",
			conf: &Configuration{
				afterFlag:      10,
				beforeFlag:     0,
				contextFlag:    0,
				countFlag:      false,
				ignoreCaseFlag: true,
				invertFlag:     false,
				fixedFlag:      false,
				lineNumFlag:    false,
				regexPattern:   "бассейн",
				files:          []string{"./test_files/grocery_list.txt"},
			},
			want: "воду для бассейна 2500000 л.\n" +
				"silva s14(черная) 1шт\n" +
				"silva s14(красная) 1шт\n" +
				"диск с Return to Versailles 1 шт\n" +
				"батл пасс 1 шт\n",
		},

		{name: "afterFlag_3_overlap",
			conf: &Configuration{
				afterFlag:      3,
				beforeFlag:     0,
				contextFlag:    0,
				countFlag:      false,
				ignoreCaseFlag: true,
				invertFlag:     false,
				fixedFlag:      false,
				lineNumFlag:    false,
				regexPattern:   "жестком|бассейн",
				files:          []string{"./test_files/grocery_list.txt"},
			},
			want: "место на жестком диске +2 ТБ\n" +
				"дверь от Chrome Hearts 1 шт.\n" +
				"воду для бассейна 2500000 л.\n" +
				"silva s14(черная) 1шт\n" +
				"silva s14(красная) 1шт\n" +
				"диск с Return to Versailles 1 шт\n",
		},

		{name: "beforeFlag_1",
			conf: &Configuration{
				afterFlag:      0,
				beforeFlag:     2,
				contextFlag:    0,
				countFlag:      false,
				ignoreCaseFlag: true,
				invertFlag:     false,
				fixedFlag:      false,
				lineNumFlag:    false,
				regexPattern:   "арбуз",
				files:          []string{"./test_files/grocery_list.txt"},
			},
			want: "авокадо 0.5 кг.\n" +
				"бензопила 0 шт.\n" +
				"АРБУЗ 4 шт.\n",
		},

		{name: "beforeFlag_2",
			conf: &Configuration{
				afterFlag:      0,
				beforeFlag:     4,
				contextFlag:    0,
				countFlag:      false,
				ignoreCaseFlag: true,
				invertFlag:     false,
				fixedFlag:      false,
				lineNumFlag:    false,
				regexPattern:   "молоко",
				files:          []string{"./test_files/grocery_list.txt"},
			},
			want: "Хлеб 2 шт.\n" +
				"молоко 2 л.\n",
		},

		{name: "beforeFlag_3_overlap",
			conf: &Configuration{
				afterFlag:      0,
				beforeFlag:     10,
				contextFlag:    0,
				countFlag:      false,
				ignoreCaseFlag: true,
				invertFlag:     false,
				fixedFlag:      false,
				lineNumFlag:    false,
				regexPattern:   "молоко|бензопила",
				files:          []string{"./test_files/grocery_list.txt"},
			},
			want: "Хлеб 2 шт.\n" +
				"молоко 2 л.\n" +
				"авокадо 0.5 кг.\n" +
				"бензопила 0 шт.\n",
		},

		{name: "contextFlag_1",
			conf: &Configuration{
				afterFlag:      0,
				beforeFlag:     0,
				contextFlag:    4,
				countFlag:      false,
				ignoreCaseFlag: true,
				invertFlag:     false,
				fixedFlag:      false,
				lineNumFlag:    false,
				regexPattern:   "финики",
				files:          []string{"./test_files/grocery_list.txt"},
			},
			want: "АРБУЗ 4 шт.\n" +
				"билеты на сектор газа 2 шт.\n" +
				"права категории С 1 шт.\n" +
				"место на кладбище 1 шт.\n" +
				"финики 7 кг.\n" +
				"конфеты 10 кг.\n" +
				"изотоп калия K40 1 г.\n" +
				"место на жестком диске +2 ТБ\n" +
				"дверь от Chrome Hearts 1 шт.\n",
		},

		{name: "contextFlag_2",
			conf: &Configuration{
				afterFlag:      0,
				beforeFlag:     0,
				contextFlag:    4,
				countFlag:      false,
				ignoreCaseFlag: true,
				invertFlag:     false,
				fixedFlag:      false,
				lineNumFlag:    false,
				regexPattern:   "молоко",
				files:          []string{"./test_files/grocery_list.txt"},
			},
			want: "Хлеб 2 шт.\n" +
				"молоко 2 л.\n" +
				"авокадо 0.5 кг.\n" +
				"бензопила 0 шт.\n" +
				"АРБУЗ 4 шт.\n" +
				"билеты на сектор газа 2 шт.\n",
		},

		{name: "contextFlag_3",
			conf: &Configuration{
				afterFlag:      0,
				beforeFlag:     0,
				contextFlag:    2,
				countFlag:      false,
				ignoreCaseFlag: true,
				invertFlag:     false,
				fixedFlag:      false,
				lineNumFlag:    false,
				regexPattern:   "батл пасс",
				files:          []string{"./test_files/grocery_list.txt"},
			},
			want: "silva s14(красная) 1шт\n" +
				"диск с Return to Versailles 1 шт\n" +
				"батл пасс 1 шт\n",
		},

		{name: "contextFlag_4",
			conf: &Configuration{
				afterFlag:      0,
				beforeFlag:     0,
				contextFlag:    2000,
				countFlag:      false,
				ignoreCaseFlag: true,
				invertFlag:     false,
				fixedFlag:      false,
				lineNumFlag:    false,
				regexPattern:   "Chrome Hearts",
				files:          []string{"./test_files/grocery_list.txt"},
			},
			want: "Хлеб 2 шт.\n" +
				"молоко 2 л.\n" +
				"авокадо 0.5 кг.\n" +
				"бензопила 0 шт.\n" +
				"АРБУЗ 4 шт.\n" +
				"билеты на сектор газа 2 шт.\n" +
				"права категории С 1 шт.\n" +
				"место на кладбище 1 шт.\n" +
				"финики 7 кг.\n" +
				"конфеты 10 кг.\n" +
				"изотоп калия K40 1 г.\n" +
				"место на жестком диске +2 ТБ\n" +
				"дверь от Chrome Hearts 1 шт.\n" +
				"воду для бассейна 2500000 л.\n" +
				"silva s14(черная) 1шт\n" +
				"silva s14(красная) 1шт\n" +
				"диск с Return to Versailles 1 шт\n" +
				"батл пасс 1 шт\n",
		},

		{name: "contextFlag_5_overlap",
			conf: &Configuration{
				afterFlag:      0,
				beforeFlag:     0,
				contextFlag:    4,
				countFlag:      false,
				ignoreCaseFlag: true,
				invertFlag:     false,
				fixedFlag:      false,
				lineNumFlag:    false,
				regexPattern:   "финики|изотоп",
				files:          []string{"./test_files/grocery_list.txt"},
			},
			want: "АРБУЗ 4 шт.\n" +
				"билеты на сектор газа 2 шт.\n" +
				"права категории С 1 шт.\n" +
				"место на кладбище 1 шт.\n" +
				"финики 7 кг.\n" +
				"конфеты 10 кг.\n" +
				"изотоп калия K40 1 г.\n" +
				"место на жестком диске +2 ТБ\n" +
				"дверь от Chrome Hearts 1 шт.\n" +
				"воду для бассейна 2500000 л.\n" +
				"silva s14(черная) 1шт\n",
		},

		{name: "line_counter_1",
			conf: &Configuration{
				afterFlag:      0,
				beforeFlag:     0,
				contextFlag:    0,
				countFlag:      true,
				ignoreCaseFlag: true,
				invertFlag:     false,
				fixedFlag:      false,
				lineNumFlag:    false,
				regexPattern:   "финики|изотоп",
				files:          []string{"./test_files/grocery_list.txt"},
			},
			want: "2\n",
		},

		{name: "line_counter_2",
			conf: &Configuration{
				afterFlag:      0,
				beforeFlag:     0,
				contextFlag:    0,
				countFlag:      true,
				ignoreCaseFlag: true,
				invertFlag:     false,
				fixedFlag:      false,
				lineNumFlag:    false,
				regexPattern:   "диск",
				files:          []string{"./test_files/grocery_list.txt"},
			},
			want: "2\n",
		},

		{name: "line_counter_3",
			conf: &Configuration{
				afterFlag:      0,
				beforeFlag:     0,
				contextFlag:    0,
				countFlag:      true,
				ignoreCaseFlag: true,
				invertFlag:     false,
				fixedFlag:      false,
				lineNumFlag:    false,
				regexPattern:   "",
				files:          []string{"./test_files/grocery_list.txt"},
			},
			want: "18\n",
		},

		{name: "invert_flag_1",
			conf: &Configuration{
				afterFlag:      0,
				beforeFlag:     0,
				contextFlag:    0,
				countFlag:      false,
				ignoreCaseFlag: true,
				invertFlag:     true,
				fixedFlag:      false,
				lineNumFlag:    false,
				regexPattern:   "",
				files:          []string{"./test_files/grocery_list.txt"},
			},
			want: "Совпадений не найдено.\n",
		},

		{name: "line_num_1",
			conf: &Configuration{
				afterFlag:      0,
				beforeFlag:     0,
				contextFlag:    0,
				countFlag:      false,
				ignoreCaseFlag: true,
				invertFlag:     false,
				fixedFlag:      false,
				lineNumFlag:    true,
				regexPattern:   "финики",
				files:          []string{"./test_files/grocery_list.txt"},
			},
			want: "9:финики 7 кг.\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			myWriterInstance.Clean()
			app(tt.conf, &myWriterInstance)
			if gotWriter := myWriterInstance.GetString(); gotWriter != tt.want {
				t.Errorf("app(): \n%v\nwant: \n%v", gotWriter, tt.want)
			}
		})
	}
}
