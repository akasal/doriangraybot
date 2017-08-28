package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"github.com/coreos/pkg/flagutil"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	s "strings"
	"github.com/hypebeast/go-osc/osc"

	"github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
	"net/http"
	"time"
	//"github.com/hypebeast/go-osc/osc"
)
var p = fmt.Println

//the following is for socketio
type Channel struct {
	Channel string `json:"channel"`
}
type Message struct {
	Id      int    `json:"id"`
	Channel string `json:"channel"`
	Text    string `json:"text"`
}


func main() {

	//socketio stuff
	server := gosocketio.NewServer(transport.GetDefaultWebsocketTransport())
	server.On(gosocketio.OnConnection, func(c *gosocketio.Channel) {
		log.Println("Connected")
		c.Emit("/message", Message{10, "main", "using emit"})
		c.Join("test")
		c.BroadcastTo("test", "/message", Message{10, "main", "using broadcast"})
//socketio continues at the end



	flags := flag.NewFlagSet("user-auth", flag.ExitOnError)
	consumerKey := flags.String("consumer-key", "OEPcDrSLMtC5pUoNNjcyewM9L", "Twitter Consumer Key")
	consumerSecret := flags.String("consumer-secret", "Lr0z3mWKbyCXTaegpzjpZh5biIx1osvdZRGf11eMUJGBVhf3Zf", "Twitter Consumer Secret")
	accessToken := flags.String("access-token", "15905105-7m7WYnV1m2hXUlusZd706GXCCDqRIUcgLj35J9SzU", "Twitter Access Token")
	accessSecret := flags.String("access-secret", "XXnJup135a0Q1LGZH3o4mB12wFyn7U9oAWjSSySBBCgiH", "Twitter Access Secret")
	flags.Parse(os.Args[1:])
	flagutil.SetFlagsFromEnv(flags, "TWITTER")

	if *consumerKey == "" || *consumerSecret == "" || *accessToken == "" || *accessSecret == "" {
		log.Fatal("Consumer key/secret and Access token/secret required")
	}

	config := oauth1.NewConfig(*consumerKey, *consumerSecret)
	token := oauth1.NewToken(*accessToken, *accessSecret)
	// OAuth1 http.Client will automatically authorize Requests
	httpClient := config.Client(oauth1.NoContext, token)

	// Twitter Client
	client := twitter.NewClient(httpClient)

	// Convenience Demux demultiplexed stream messages
	demux := twitter.NewSwitchDemux()

	demux.Tweet = func(tweet *twitter.Tweet) {
		//p("Contains:  ", s.Contains("test", "es"))
		//p("Contains:  ", s.Contains(tweet.Text, "love"))
		//if (s.Contains(tweet.Text, "love"))||(s.Contains(tweet.Text, "hate"))||(s.Contains(tweet.Text, "trump")){
		if (s.Contains(tweet.Text, "love"))||(s.Contains(tweet.Text, "hate"))||(s.Contains(tweet.Text, "LOVE"))||(s.Contains(tweet.Text, "HATE"))||(s.Contains(tweet.Text, "Love"))||(s.Contains(tweet.Text, "Hate")){
			fmt.Println(tweet.Text)
			client := osc.NewClient("localhost", 8765)
			msg := osc.NewMessage("/dorian/address")
			//msg.Append(int32(111))
			//msg.Append(true)
			msg.Append(tweet.Text) //was "facechange"
			client.Send(msg)
			c.Emit("/message", Message{10, "main", tweet.Text})
		}
		//fmt.Println(tweet.Text)
	}






	demux.DM = func(dm *twitter.DirectMessage) {
		//fmt.Println(dm.SenderID)
	}
	demux.Event = func(event *twitter.Event) {
		//fmt.Printf("%#v\n", event)
	}

	fmt.Println(" ")
	fmt.Println("Starting Stream for Dilfer<3...")
	fmt.Println(" ")
	fmt.Println("This includes all tweets that include the words love OR hate.")
	fmt.Println(" ")

	// FILTER
	//filterParams := &twitter.StreamFilterParams{
	//	Track:         []string{"love"},
	//	StallWarnings: twitter.Bool(true),
	//}
	//stream, err := client.Streams.Filter(filterParams)
	//if err != nil {
	//	log.Fatal(err)
	//}

	// USER (quick test: auth'd user likes a tweet -> event)
	// userParams := &twitter.StreamUserParams{
	// 	StallWarnings: twitter.Bool(true),
	// 	With:          "followings",
	// 	Language:      []string{"en"},
	// }
	// stream, err := client.Streams.User(userParams)
	// if err != nil {
	// 	log.Fatal(err)
	// }



	// SAMPLE
	 sampleParams := &twitter.StreamSampleParams{
	 	StallWarnings: twitter.Bool(true),
	 }
	 stream, err := client.Streams.Sample(sampleParams)
	 if err != nil {
	 	log.Fatal(err)
	 }

	// Receive messages until stopped or stream quits
	go demux.HandleChan(stream.Messages)

	// Wait for SIGINT and SIGTERM (HIT CTRL-C)
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)

	fmt.Println("Stopping Stream for Alpaca_xxx...")
	stream.Stop()








//backto closing socketio
	})
	server.On(gosocketio.OnDisconnection, func(c *gosocketio.Channel) {
		log.Println("Disconnected")
	})
	server.On("/join", func(c *gosocketio.Channel, channel Channel) string {
		time.Sleep(2 * time.Second)
		log.Println("Client joined to ", channel.Channel)
		return "joined to " + channel.Channel
	})
	serveMux := http.NewServeMux()
	serveMux.Handle("/socket.io/", server)
	log.Println("Starting server...")
	log.Panic(http.ListenAndServe(":3811", serveMux))
	//end socketio



	}
