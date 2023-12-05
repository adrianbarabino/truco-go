package main

import (
    "strconv"
    "strings"
    // "database/sql"
    "encoding/json"
    "fmt"
    "math/rand"
    //"io/ioutil"
    "log"
    "net/http"
	"html/template"
    // _ "github.com/go-sql-driver/mysql"
     "github.com/gorilla/mux"
     "github.com/gorilla/websocket"

)

type Card struct {
    ID      int    `json:"id"`
    Number   int `json:"number"`
    Type string `json:"type"`
}
var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}
type Game struct {
    Players     []*websocket.Conn
    PlayerTurn  int
    RoundCards []Card // Cartas jugadas en la ronda actual

    // Otros estados del juego como el tablero, cartas, etc.
}
func determinarGanador(card1, card2 Card) (ganador int, puntos int) {
    card1Value := getCardValue(card1)
    card2Value := getCardValue(card2)

    if card1Value > card2Value {
        return 0, card1Value // Suponiendo que el jugador 0 jugó la card1
    } else if card2Value > card1Value {
        return 1, card2Value // Suponiendo que el jugador 1 jugó la card2
    }
    return -1, 0 // Empate, ningún jugador gana puntos
}

func getCardValue(card Card) int {
    // Definir el valor especial de algunas cartas en el truco
    specialCards := map[string]int{
        "1-espada": 14,
        "1-basto": 13,
        "7-espada": 12,
        "7-oro": 11,
        "3-copa": 10,
        "3-oro": 10,
        "3-espada": 10,
        "3-basto": 10,
        "2-copa": 9,
        "2-oro": 9,
        "2-espada": 9,
        "2-basto": 9,
        "1-copa": 8,
        "1-oro": 8,
        "12-espada": 7,
        "12-basto": 7,
        "12-oro": 7,
        "12-copa": 7,
        "11-espada": 6,
        "11-basto": 6,
        "11-oro": 6,
        "11-copa": 6,
        "10-espada": 5,
        "10-basto": 5,
        "10-oro": 5,
        "10-copa": 5,
        "7-copa": 4,
        "7-basto": 4,
        "6-espada": 3,
        "6-basto": 3,
        "6-oro": 3,
        "6-copa": 3,
        "5-espada": 2,
        "5-basto": 2,
        "5-oro": 2,
        "5-copa": 2,
        "4-espada": 1,
        "4-basto": 1,
        "4-oro": 1,
        "4-copa": 1,
    }

    cardKey := fmt.Sprintf("%d-%s", card.Number, card.Type)
    value, exists := specialCards[cardKey]
    if exists {
        return value
    }
    return card.Number // Valor por defecto si no es una carta especial
}

// Función para enviar un mensaje a todos los jugadores
func (g *Game) broadcastMessage(message string) {
    for _, player := range g.Players {
        if player != nil {
            player.WriteMessage(websocket.TextMessage, []byte(message))
        }
    }
}

// Función para manejar el turno de los jugadores
func (g *Game) nextTurn() {
    g.PlayerTurn = (g.PlayerTurn + 1) % len(g.Players)
    g.broadcastMessage(fmt.Sprintf("Turno del jugador %d", g.PlayerTurn+1))
    // enviale solo al jugador que le toca el turno, el permiso para jugar
    g.Players[g.PlayerTurn].WriteMessage(websocket.TextMessage, []byte("Tu turno"))

}

// Función para manejar la conexión WebSocket y las acciones de los jugadores
func handlePlayerActions(w http.ResponseWriter, r *http.Request) {
    // cuando se ejecuta esta funcion handlePlayerActions?
    // cuando se conecta un jugador a la ruta /ws


    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println(err)
        return
    }

    // Añadir el jugador al juego
    newGame.Players = append(newGame.Players, conn)

    // Comprobar si todos los jugadores están conectados para empezar
    if len(newGame.Players) == 2 {
        newGame.broadcastMessage("Todos los jugadores están conectados. Comienza el juego.")
        newGame.nextTurn() // Inicia el turno del primer jugador
    }

    for {
        // Leer mensaje del jugador
        _, message, err := conn.ReadMessage()
        if err != nil {
            log.Println("Error al leer mensaje:", err)
            break
        }

        // Procesar acción del jugador
        log.Printf("Acción del jugador: %s", message)
        // Aquí procesas la lógica del juego dependiendo de la acción

        // Pasar al siguiente turno
        newGame.nextTurn()
    }
}


// Nuevo juego
var newGame = Game{
    Players: make([]*websocket.Conn, 0),
}
type Player struct {
    ID      int    `json:"id"`
    Name   string `json:"name"`
}

type Hand struct {
    ID   int    `json:"id"`
    Cards []Card `json:"cards"`
}



var plantillas = template.Must(template.ParseFiles("carta.html"));
func cargarPlantilla(w http.ResponseWriter, nombre_plantilla string, pagina int){
	plantillas.ExecuteTemplate(w, nombre_plantilla + ".html", 1);
}
func main() {
    // // Leer la configuración de data.json
    // configFile, err := ioutil.ReadFile("data.json")
    // if err != nil {
    //     log.Fatal("Error al leer archivo de configuración: ", err)
    // }
	var b = 0;
	var c = 0;
	var types []string
	var allCards []Card
	types = append(types, "basto");

	types = append(types, "oro");

	types = append(types, "copa");

	types = append(types, "espada");

	for i := 0; i < 40; i++ {
		b++;
		if(b == 8){
			b = 10;
		}

		if(b == 9){
			b = 11;
		}

		if(b == 10){
			b = 12;
		}
		var newCard Card
		newCard.ID = i;
		newCard.Number = b;
		newCard.Type = types[c];
		if(b == 12){
			b = 0;
			if(c < 3){
				c++;
			}
		}
		
		allCards = append(allCards, newCard)
		
	}
    // var config DBConfig
    // err = json.Unmarshal(configFile, &config)
    // if err != nil {
    //     log.Fatal("Error al parsear archivo de configuración: ", err)
    // }

    // // Conectar a la base de datos
    // db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", config.User, config.Password, config.Host, config.Port, config.DBName))
    // if err != nil {
    //     log.Fatal(err)
    // }
    // defer db.Close()

     router := mux.NewRouter()

    // Resto del código...
    


	

    router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        
		cargarPlantilla(w, "carta", 1);

    })


    router.HandleFunc("/cartas", func(w http.ResponseWriter, r *http.Request) {
        
        w.Header().Set("Content-Type", "application/json")
        w.Header().Set("Access-Control-Allow-Origin", "*")
        
        w.Header().Set("Origin", "*")
        json.NewEncoder(w).Encode(allCards)
    })
// Manejador para la conexión WebSocket
router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
   

    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println(err)
        return
    }

    // Añadir el jugador al juego
    newGame.Players = append(newGame.Players, conn)

    // Comprobar si todos los jugadores están conectados para empezar
    if len(newGame.Players) == 2 {
        newGame.broadcastMessage("Todos los jugadores están conectados. Comienza el juego.")
        
        // Repartir cartas a los jugadores al azar y sin repetidos
        var playerCards []Hand

// Hacer una copia del slice allCards para mantener los índices consistentes
availableCards := make([]Card, len(allCards))
copy(availableCards, allCards)

for i := 0; i < 2; i++ {
    var newHand Hand
    newHand.ID = i
    newHand.Cards = make([]Card, 0)

    for j := 0; j < 3; j++ {
        // Generar un índice aleatorio para la carta
        randomIndex := rand.Intn(len(availableCards))
        randomCard := availableCards[randomIndex]

        // Agregar la carta seleccionada a la nueva mano
        newHand.Cards = append(newHand.Cards, randomCard)

        // Eliminar la carta seleccionada de availableCards
        availableCards = append(availableCards[:randomIndex], availableCards[randomIndex+1:]...)
    }

    playerCards = append(playerCards, newHand)
}



        // generar el json de las cartas de cada jugador
        jsonCards1, _ := json.Marshal(playerCards[0])
        jsonCards2, _ := json.Marshal(playerCards[1])

        // enviar las cartas a cada jugador en formato json y con un indicador en el mensaje para diferenciarlo
        newGame.Players[0].WriteMessage(websocket.TextMessage, []byte("Cartas: " + string(jsonCards1)))
        newGame.Players[1].WriteMessage(websocket.TextMessage, []byte("Cartas: " + string(jsonCards2)))

        newGame.nextTurn() // Inicia el turno del primer jugador

        // enviar solo las cartas de cada jugador a cada jugador
        // no enviar el json completo 
    }

    for {
        // Leer mensaje del jugador
        _, message, err := conn.ReadMessage()
        if err != nil {
            log.Println("Error al leer mensaje:", err)
            break
        }

        // Procesar acción del jugador
        log.Printf("Acción del jugador: %s", message)
        // si el message es "Carta movida: carta-19" extraer el numero del final y enviarlo a todos los jugadores



        if strings.HasPrefix(string(message), "Carta movida:") {
            // Extrae el substring que sigue después de "Carta movida: carta-"
            cardIDStr := strings.TrimSpace(strings.TrimPrefix(string(message), "Carta movida: carta-"))
            
            // Convierte el substring a un entero
            cardID, err := strconv.Atoi(cardIDStr)
            if err != nil {
                log.Println("Error al convertir la ID de la carta:", err)
                return
            }

            log.Println("ID de la carta movida:", cardID)
    // Encuentra la carta por ID y agrégala a RoundCards
    for _, card := range allCards { // Asume que allCards contiene todas las cartas
        if card.ID == cardID {
            log.Println("Carta movida:", card)
            log.Println("Carta actual:", card.ID);
            log.Println("Carta actual:", card.Number);
            log.Println("Carta actual:", card.Type);
            
            newGame.RoundCards = append(newGame.RoundCards, card)
            break
        }
    }

    // Envía la carta movida al otro jugador
    if newGame.PlayerTurn == 0 {
        newGame.Players[1].WriteMessage(websocket.TextMessage, []byte("Carta movida: " + string(message[14:])))
        } else {
            newGame.Players[0].WriteMessage(websocket.TextMessage, []byte("Carta movida: " + string(message[14:])))
        }

    // Si se han jugado 2 cartas, determina el ganador de la ronda
    if len(newGame.RoundCards) == 2 {
        ganador, puntos := determinarGanador(newGame.RoundCards[0], newGame.RoundCards[1])
        // Lógica para actualizar el estado del juego con el ganador y los puntos
        // ...
        newGame.broadcastMessage("El ganador es: " + strconv.Itoa(ganador+1) + " con " + strconv.Itoa(puntos) + " puntos")


        newGame.RoundCards = nil // Limpia las cartas para la siguiente ronda
    }

    newGame.nextTurn()

}




        // Aquí procesas la lógica del juego dependiendo de la acción

        // Pasar al siguiente turno
    }
})

    router.PathPrefix("/").Handler(http.FileServer(http.Dir("./assets/")))


    fmt.Println("Server is running at :80")
    log.Fatal(http.ListenAndServe(":80", router))
}


