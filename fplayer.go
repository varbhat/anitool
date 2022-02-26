package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/aerogo/http/client"
)

// class FPlayer(private val getSize:Boolean): Extractor() {
//     override fun getStreamLinks(name: String, url: String): Episode.StreamLinks {
//         val apiLink = url.replace("/v/","/api/source/")
//         val tempQuality = mutableListOf<Episode.Quality>()
//         try{
//         val jsonResponse = Json.decodeFromString<JsonObject>(Jsoup.connect(apiLink).ignoreContentType(true)
//             .header("referer",url)
//             .post().body().text())

//         if(jsonResponse["success"].toString() == "true") {
//             val a = arrayListOf<Deferred<*>>()
//             runBlocking {
//                 jsonResponse.jsonObject["data"]!!.jsonArray.forEach {
//                     a.add(async {
//                         tempQuality.add(
//                             Episode.Quality(
//                                 it.jsonObject["file"].toString().trim('"'),
//                                 it.jsonObject["label"].toString().trim('"'),
//                                 if(getSize) getSize(it.jsonObject["file"].toString().trim('"')) else null
//                             )
//                         )
//                     })
//                 }
//             }
//         }
//         }catch (e:Exception){ toastString(e.toString()) }
//         return Episode.StreamLinks(
//             name,
//             tempQuality,
//             null
//         )
//     }

// }

type FplayerResp struct {
	Data []Link `json:"data"`
}

func Fplayer(iurl string) []Link {
	apiurl := strings.Replace(iurl, "/v/", "/api/source/", -1)
	response, err := client.Post(apiurl).End()
	if err != nil {
		log.Fatal(err)

	}
	var Fresp FplayerResp
	if err := json.Unmarshal(response.Bytes(), &Fresp); err != nil {
		return []Link{}
	}

	for _, eachSource := range Fresp.Data {
		fmt.Println("File ", eachSource.File)
		fmt.Println("Label ", eachSource.Label)
		fmt.Println("Type ", eachSource.Type)
	}

	return []Link{}
}
