package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"mangahub/internal/auth"
)

type seedManga struct {
	Title        string   `json:"title"`
	Author       string   `json:"author"`
	Artist       string   `json:"artist"`
	Genres       []string `json:"genres"`
	Status       string   `json:"status"`
	ChapterCount int      `json:"chapter_count"`
	VolumeCount  int      `json:"volume_count"`
	Description  string   `json:"description"`
	CoverURL     string   `json:"cover_url"`
	Year         int      `json:"year"`
	Rating       float64  `json:"rating"`
	Rank         int      `json:"rank"`
}

func seedData(db *sql.DB) error {
	log.Println("[seed] seeding manga...")

	data, err := os.ReadFile("data/seeds/manga.json")
	if err != nil {
		log.Printf("[seed] no manga.json found, using built-in seed")
		data = builtinMangaSeed
	}

	var manga []seedManga
	if err := json.Unmarshal(data, &manga); err != nil {
		return err
	}

	for _, m := range manga {
		id := uuid.New()
		_, err := db.Exec(`
			INSERT INTO manga (id,title,author,artist,genres,status,chapter_count,volume_count,
			                   description,cover_url,year,rating,popularity_rank)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)
			ON CONFLICT (title) DO UPDATE 
			SET cover_url = EXCLUDED.cover_url,
			    rating = EXCLUDED.rating,
			    chapter_count = EXCLUDED.chapter_count`,
			id, m.Title, m.Author, m.Artist,
			pq.Array(m.Genres), m.Status,
			m.ChapterCount, m.VolumeCount,
			m.Description, m.CoverURL,
			m.Year, m.Rating, m.Rank,
		)
		if err != nil {
			log.Printf("[seed] manga %s: %v", m.Title, err)
		} else {
			log.Printf("[seed] inserted manga: %s", m.Title)
		}
	}

	// Create admin user
	log.Println("[seed] creating admin user...")
	hash, _ := auth.HashPassword("admin1234")
	adminID := uuid.New()
	_, err = db.Exec(`
		INSERT INTO users (id,username,email,password_hash,role)
		VALUES ($1,'admin','admin@mangahub.dev',$2,'admin')
		ON CONFLICT (email) DO NOTHING`, adminID, hash)
	if err != nil {
		log.Printf("[seed] admin user: %v", err)
	}
	_, _ = db.Exec(`INSERT INTO user_settings (user_id) VALUES ($1) ON CONFLICT DO NOTHING`, adminID)

	// Create demo user
	log.Println("[seed] creating demo user...")
	demoHash, _ := auth.HashPassword("demo1234")
	demoID := uuid.New()
	_, err = db.Exec(`
		INSERT INTO users (id,username,email,password_hash,role)
		VALUES ($1,'demouser','demo@mangahub.dev',$2,'user')
		ON CONFLICT (email) DO NOTHING`, demoID, demoHash)
	if err != nil {
		log.Printf("[seed] demo user: %v", err)
	}
	_, _ = db.Exec(`INSERT INTO user_settings (user_id) VALUES ($1) ON CONFLICT DO NOTHING`, demoID)

	log.Println("[seed] done ✓")
	return nil
}

var builtinMangaSeed = []byte(`[
  {"title":"One Piece","author":"Eiichiro Oda","artist":"Eiichiro Oda","genres":["Action","Adventure","Comedy","Fantasy","Shounen"],"status":"ongoing","chapter_count":1105,"volume_count":107,"description":"Monkey D. Luffy sets off on a journey to find the legendary treasure One Piece and become King of the Pirates.","cover_url":"https://uploads.mangadex.org/covers/a1c02e1c-5d6e-445a-8b83-a75d31599a09/5886994a-3841-4770-9883-fc148a045763.jpg","year":1997,"rating":9.2,"rank":1},
  {"title":"Berserk","author":"Kentaro Miura","artist":"Kentaro Miura","genres":["Action","Adventure","Dark Fantasy","Horror","Seinen"],"status":"ongoing","chapter_count":374,"volume_count":41,"description":"The story follows Guts, a lone mercenary warrior, and his tumultuous relationship with Griffith, the leader of a mercenary band.","cover_url":"https://uploads.mangadex.org/covers/8015c324-a27a-4691-b45c-bc03485777bd/6d7cc652-32b0-449d-9293-61a03975005b.jpg","year":1989,"rating":9.4,"rank":2},
  {"title":"Vinland Saga","author":"Makoto Yukimura","artist":"Makoto Yukimura","genres":["Action","Adventure","Drama","Historical","Seinen"],"status":"ongoing","chapter_count":196,"volume_count":27,"description":"Set in medieval Europe during the Danish invasion of England, Thorfinn seeks to avenge his father's death.","cover_url":"https://uploads.mangadex.org/covers/337fb31a-e525-4b31-8f35-9f5b6678b668/ce7b9f84-f3d3-4610-8547-064c5d6ba81e.jpg","year":2005,"rating":9.1,"rank":3},
  {"title":"Vagabond","author":"Takehiko Inoue","artist":"Takehiko Inoue","genres":["Action","Adventure","Drama","Historical","Seinen"],"status":"hiatus","chapter_count":327,"volume_count":37,"description":"A fictionalized retelling of the life of Miyamoto Musashi, the greatest samurai who ever lived.","cover_url":"https://uploads.mangadex.org/covers/d8a3916b-19e3-470a-810d-374/337fb31a-e525-4b31-8f35-9f5b6678b668.jpg","year":1998,"rating":9.2,"rank":4},
  {"title":"Fullmetal Alchemist","author":"Hiromu Arakawa","artist":"Hiromu Arakawa","genres":["Action","Adventure","Fantasy","Military","Shounen"],"status":"completed","chapter_count":116,"volume_count":27,"description":"Two brothers use alchemy and search for the Philosopher's Stone to restore their bodies after a failed transmutation.","cover_url":"https://uploads.mangadex.org/covers/8efcb55c-1961-469d-8f24-9b8e84323e07/d7e0081d-e6b0-4966-a496-e30a597a4773.jpg","year":2001,"rating":9.0,"rank":5},
  {"title":"Attack on Titan","author":"Hajime Isayama","artist":"Hajime Isayama","genres":["Action","Drama","Fantasy","Horror","Shounen","Thriller"],"status":"completed","chapter_count":139,"volume_count":34,"description":"Humanity lives behind walls to protect themselves from titans. Eren Yeager vows to kill all titans after his mother is eaten.","cover_url":"https://uploads.mangadex.org/covers/30468241-d610-4497-a725-b8253a63ec50/23d130ed-3679-455b-866-9a2eb51833d7.jpg","year":2009,"rating":8.8,"rank":6},
  {"title":"Demon Slayer","author":"Koyoharu Gotouge","artist":"Koyoharu Gotouge","genres":["Action","Fantasy","Historical","Shounen","Supernatural"],"status":"completed","chapter_count":205,"volume_count":23,"description":"Tanjiro Kamado sets out to avenge his slaughtered family and cure his sister, who has been turned into a demon.","cover_url":"https://uploads.mangadex.org/covers/d45903e1-7006-444b-b13c-cf80efcb7a9c/d2825b29-e0c4-42b7-a3f2-10f84501b446.jpg","year":2016,"rating":8.7,"rank":7},
  {"title":"Naruto","author":"Masashi Kishimoto","artist":"Masashi Kishimoto","genres":["Action","Adventure","Comedy","Martial Arts","Shounen"],"status":"completed","chapter_count":700,"volume_count":72,"description":"Naruto Uzumaki, a young ninja, seeks recognition from his peers and dreams of becoming the Hokage.","cover_url":"https://uploads.mangadex.org/covers/466085a6-0f8f-4ed1-8f5a-714a6e35328e/6573e049-7c15-468e-976d-d779a528e5c1.jpg","year":1999,"rating":8.1,"rank":8},
  {"title":"Dragon Ball","author":"Akira Toriyama","artist":"Akira Toriyama","genres":["Action","Adventure","Comedy","Fantasy","Martial Arts","Shounen"],"status":"completed","chapter_count":519,"volume_count":42,"description":"Son Goku goes on a quest to find the seven Dragon Balls that can grant any wish when gathered.","cover_url":"https://uploads.mangadex.org/covers/f5fb2b23-1d01-4475-8461-9e8c0b561c21/04fa05a2-369f-4315-992a-e8d1979b003a.jpg","year":1984,"rating":8.2,"rank":9},
  {"title":"Hunter x Hunter","author":"Yoshihiro Togashi","artist":"Yoshihiro Togashi","genres":["Action","Adventure","Fantasy","Shounen","Supernatural"],"status":"hiatus","chapter_count":400,"volume_count":37,"description":"Gon Freecss discovers his father is a famous hunter and sets out to become one himself while making many friends.","cover_url":"https://uploads.mangadex.org/covers/f3c55985-8025-4c07-88C4-d92a95535359/11b330ec-4f6d-4959-adb5-3738b312e5ec.jpg","year":1998,"rating":9.0,"rank":10},
  {"title":"Death Note","author":"Tsugumi Ohba","artist":"Takeshi Obata","genres":["Drama","Mystery","Psychological","Supernatural","Thriller"],"status":"completed","chapter_count":108,"volume_count":12,"description":"Light Yagami finds the Death Note and uses it to rid the world of criminals, leading to a cat-and-mouse game with detective L.","cover_url":"https://uploads.mangadex.org/covers/66396e95-7098-4c02-9a4e-bb9936d94dda/762cded0-28ba-4f7f-abc9-698f2abb6295.jpg","year":2003,"rating":8.7,"rank":11},
  {"title":"Tokyo Ghoul","author":"Sui Ishida","artist":"Sui Ishida","genres":["Action","Dark Fantasy","Horror","Seinen","Supernatural","Tragedy"],"status":"completed","chapter_count":179,"volume_count":16,"description":"Ken Kaneki survives a deadly encounter with a ghoul and becomes a half-ghoul, half-human hybrid.","cover_url":"https://uploads.mangadex.org/covers/682cf51b-10f8-43d9-ab5f-0b96a8c6766f/e2dab7d0-d5a8-4e8d-b274-069ee73D35e2.jpg","year":2011,"rating":8.2,"rank":12},
  {"title":"My Hero Academia","author":"Kohei Horikoshi","artist":"Kohei Horikoshi","genres":["Action","Comedy","School","Shounen","Super Power"],"status":"completed","chapter_count":430,"volume_count":41,"description":"In a world where most people have superpowers, Izuku Midoriya is born without powers but still dreams of becoming a hero.","cover_url":"https://uploads.mangadex.org/covers/d3896504-03a1-432a-b07e-ebbe2a826649/842d307d-b2e5-4e7f-8547-0b0b042f377e.jpg","year":2014,"rating":8.2,"rank":13},
  {"title":"Chainsaw Man","author":"Tatsuki Fujimoto","artist":"Tatsuki Fujimoto","genres":["Action","Dark Fantasy","Horror","Shounen","Supernatural"],"status":"ongoing","chapter_count":175,"volume_count":17,"description":"Denji, a young man in debt, merges with his pet chainsaw devil Pochita and becomes Chainsaw Man, a devil hunter.","cover_url":"https://uploads.mangadex.org/covers/a7774ce4-ef20-4236-b7e7-5ba5b0a9a7a6/8d7c4d9e-c2c9-4c5e-855a-baba2bdd390e.jpg","year":2018,"rating":8.8,"rank":14},
  {"title":"Spy x Family","author":"Tatsuya Endo","artist":"Tatsuya Endo","genres":["Action","Comedy","Romance","Seinen","Slice of Life"],"status":"ongoing","chapter_count":95,"volume_count":12,"description":"A spy must form a fake family to complete his mission, unknowingly adopting a telepath and marrying an assassin.","cover_url":"https://uploads.mangadex.org/covers/dab57887-4e31-4643-85af-e46107797527/157d6052-54ad-470a-b9b4-406085a8b3b5.jpg","year":2019,"rating":8.6,"rank":15},
  {"title":"Blue Period","author":"Tsubasa Yamaguchi","artist":"Tsubasa Yamaguchi","genres":["Drama","School","Seinen","Slice of Life"],"status":"ongoing","chapter_count":78,"volume_count":15,"description":"A highly academic high schooler discovers the world of art and decides to apply to Tokyo University of the Arts.","cover_url":"https://uploads.mangadex.org/covers/79707324-402b-47e2-88f5-46747df3750a/85e197d1-e945-4df3-a15e-4366ed625805.jpg","year":2017,"rating":8.7,"rank":16},
  {"title":"Mushishi","author":"Yuki Urushibara","artist":"Yuki Urushibara","genres":["Adventure","Fantasy","Historical","Mystery","Seinen","Slice of Life"],"status":"completed","chapter_count":50,"volume_count":10,"description":"Ginko travels across Japan to research Mushi, supernatural beings that exist outside of normal life.","cover_url":"https://uploads.mangadex.org/covers/6f4ef808-16cb-4034-86a0-5c6563603d2b/08c02824-7341-4566-a67b-1d7470659798.jpg","year":1999,"rating":8.7,"rank":17},
  {"title":"Planetes","author":"Makoto Yukimura","artist":"Makoto Yukimura","genres":["Drama","Romance","Sci-Fi","Seinen","Space"],"status":"completed","chapter_count":26,"volume_count":4,"description":"In the 2070s, a crew of debris collectors work in outer space to clean up orbital debris.","cover_url":"https://uploads.mangadex.org/covers/9712ed21-0be8-4228-b99b-4e14f9d2d001/a04107e3-0544-4824-b153-27cc1f7481ba.jpg","year":1999,"rating":8.6,"rank":18},
  {"title":"20th Century Boys","author":"Naoki Urasawa","artist":"Naoki Urasawa","genres":["Drama","Mystery","Sci-Fi","Seinen","Thriller"],"status":"completed","chapter_count":249,"volume_count":22,"description":"A group of childhood friends face a mysterious villain who is fulfilling the prophecies of a story they wrote as children.","cover_url":"https://uploads.mangadex.org/covers/885a0601-0ba5-4191-8848-1249b6d396a8/69591410-b5cc-447a-8f55-6b586118b6ec.jpg","year":1999,"rating":8.8,"rank":19},
  {"title":"Monster","author":"Naoki Urasawa","artist":"Naoki Urasawa","genres":["Drama","Mystery","Psychological","Seinen","Thriller"],"status":"completed","chapter_count":162,"volume_count":18,"description":"A brilliant surgeon saves a young boy's life, setting off a chain of events that leads him on a hunt for a serial killer.","cover_url":"https://uploads.mangadex.org/covers/0e28f328-86d7-400b-ba65-9850550f7572/c859ce87-738b-4afb-84a1-77884d5de315.jpg","year":1994,"rating":9.1,"rank":20},
  {"title":"JoJo's Bizarre Adventure","author":"Hirohiko Araki","artist":"Hirohiko Araki","genres":["Action","Adventure","Supernatural","Shounen","Mystery"],"status":"ongoing","chapter_count":959,"volume_count":131,"description":"The multigenerational tale of the Joestar family and their battles against supernatural forces.","cover_url":"https://uploads.mangadex.org/covers/db06173a-4933-4395-934c-6f81ea9e42e5/67566144-8395-4655-9037-33959b854371.jpg","year":1986,"rating":8.7,"rank":21},
  {"title":"Made in Abyss","author":"Akihito Tsukushi","artist":"Akihito Tsukushi","genres":["Adventure","Dark Fantasy","Mystery","Sci-Fi","Seinen"],"status":"ongoing","chapter_count":70,"volume_count":11,"description":"In a world with a mysterious giant chasm, a girl and a robot boy descend into the Abyss in search of her mother.","cover_url":"https://uploads.mangadex.org/covers/98939c36-056a-4a0a-90e9-889a7aa2e59e/46663f74-3759-4da5-b463-8848f1c22934.jpg","year":2012,"rating":8.8,"rank":31},
  {"title":"Solo Leveling","author":"Chugong","artist":"DUBU","genres":["Action","Adventure","Fantasy","Manhwa","Supernatural"],"status":"completed","chapter_count":179,"volume_count":8,"description":"In a world where hunters fight monsters, the weakest hunter, Sung Jin-Woo, gains the ability to level up in strength.","cover_url":"https://uploads.mangadex.org/covers/32ad4fc1-feaf-4bd2-bad1-5b9d1c5c29da/34537e90-6435-4309-8581-3e04e5f5a1c4.jpg","year":2018,"rating":8.5,"rank":32},
  {"title":"Black Clover","author":"Yuki Tabata","artist":"Yuki Tabata","genres":["Action","Adventure","Comedy","Fantasy","Magic","Shounen"],"status":"ongoing","chapter_count":372,"volume_count":36,"description":"Asta, a boy born without magic in a magic-filled world, dreams of becoming the Magic Emperor.","cover_url":"https://uploads.mangadex.org/covers/98394473-ba5d-4e0c-b302-36eb6e379201/22e5efcb-cb7f-441a-9eea-6677efa9d1ed.jpg","year":2015,"rating":7.9,"rank":33},
  {"title":"Fire Force","author":"Atsushi Ohkubo","artist":"Atsushi Ohkubo","genres":["Action","Fantasy","Sci-Fi","Shounen","Supernatural"],"status":"completed","chapter_count":304,"volume_count":34,"description":"In a world plagued by spontaneous human combustion, Shinra Kusakabe joins a special firefighting force.","cover_url":"https://uploads.mangadex.org/covers/8efcb55c-1961-469d-8f24-9b8e84323e07/d7e0081d-e6b0-4966-a496-e30a597a4773.jpg","year":2015,"rating":8.1,"rank":34},
  {"title":"Dungeon Meshi","author":"Ryoko Kui","artist":"Ryoko Kui","genres":["Adventure","Comedy","Fantasy","Seinen"],"status":"completed","chapter_count":97,"volume_count":14,"description":"Laios and his party explore a dungeon and survive by cooking and eating the monsters they defeat.","cover_url":"https://uploads.mangadex.org/covers/337fb31a-e525-4b31-8f35-9f5b6678b668/ce7b9f84-f3d3-4610-8547-064c5d6ba81e.jpg","year":2014,"rating":8.8,"rank":35}
]`)
