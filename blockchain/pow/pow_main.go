package pow

// 1. PoW , serve TCP ( Using netcat )
// func pow_main() {
// 	err := godotenv.Load()
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	pow.BCserver = make(chan []pow.Block)
// 	pow.Init()

// 	// start TCP and serve TCP server
// 	server, err := net.Listen("tcp", "localhost:"+os.Getenv("ADDR"))
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	defer server.Close()

// 	// 요청이 들어올때마다 새로운 Conn을 생성해야한다
// 	for {
// 		conn, err := server.Accept()
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		go pow.Conn(conn)
// 	}

// }
