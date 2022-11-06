<template>
  <div>
    <div class="dashboard-container">
      <div class="menu">
        <menu-item v-for="item in menu" :name="item.path" :count="item.count" :date="item.date"></menu-item>
      </div>
      <div class="logs">
        <div v-if="!this.connection" class="wait">
          <div class="icon">
            <font-awesome-icon icon="fa-solid fa-water" />
          </div>
          <div class="text">
            Select channel
          </div>
        </div>
        <div v-else>
          <log-item v-for="log in logs" :date="log.date" :data="log.data"></log-item>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import MenuItem from "./chunks/menu-item";
import LogItem from "./chunks/log-item";
import axios from "axios";
export default {
  components: {MenuItem, LogItem},
  data: function () {
    return {
      connection: null,
      menu: [],
      logs: [
        {
          date: "21:33 04.09.2022",
          data: {
            test: "test"
          }
        },
        {
          date: "21:33 04.09.2022",
          data: {
            test: "test"
          }
        },
        {
          date: "21:33 04.09.2022",
          data: {
            test: "test"
          }
        },
        {
          date: "21:33 04.09.2022",
          data: {
            test: "test"
          }
        },
        {
          date: "21:33 04.09.2022",
          data: {
            test: "test"
          }
        },
        {
          date: "21:33 04.09.2022",
          data: {
            test: "test"
          }
        },
        {
          date: "21:33 04.09.2022",
          data: {
            test: "test"
          }
        },
        {
          date: "21:33 04.09.2022",
          data: {
            test: "test"
          }
        },
        {
          date: "21:33 04.09.2022",
          data: {
            test: "test"
          }
        },
        {
          date: "21:33 04.09.2022",
          data: {
            test: "test"
          }
        },
        {
          date: "21:33 04.09.2022",
          data: {
            test: "test"
          }
        },
        {
          date: "21:33 04.09.2022",
          data: {
            total: 25,
            limit: 10,
            skip: 0,
            links: {
              previous: undefined,
              next: function () {},
            },
            data: [
              {
                id: '5968fcad629fa84ab65a5247',
                firstname: 'Ada',
                lastname: 'Lovelace',
                awards: null,
                known: [
                  'mathematics',
                  'computing'
                ],
                position: {
                  lat: 44.563836,
                  lng: 6.495139
                },
                description: `Augusta Ada King, Countess of Lovelace (née Byron; 10 December 1815 – 27 November 1852) was an English mathematician and writer,
            chiefly known for her work on Charles Babbage's proposed mechanical general-purpose computer,
            the Analytical Engine. She was the first to recognise that the machine had applications beyond pure calculation,
            and published the first algorithm intended to be carried out by such a machine.
            As a result, she is sometimes regarded as the first to recognise the full potential of a "computing machine" and the first computer programmer.`,
                bornAt: '1815-12-10T00:00:00.000Z',
                diedAt: '1852-11-27T00:00:00.000Z'
              }, {
                id: '5968fcad629fa84ab65a5246',
                firstname: 'Grace',
                lastname: 'Hopper',
                awards: [
                  'Defense Distinguished Service Medal',
                  'Legion of Merit',
                  'Meritorious Service Medal',
                  'American Campaign Medal',
                  'World War II Victory Medal',
                  'National Defense Service Medal',
                  'Armed Forces Reserve Medal',
                  'Naval Reserve Medal',
                  'Presidential Medal of Freedom'
                ],
                known: null,
                position: {
                  lat: 43.614624,
                  lng: 3.879995
                },
                description: `Grace Brewster Murray Hopper (née Murray; December 9, 1906 – January 1, 1992)
            was an American computer scientist and United States Navy rear admiral.
            One of the first programmers of the Harvard Mark I computer,
            she was a pioneer of computer programming who invented one of the first compiler related tools.
            She popularized the idea of machine-independent programming languages, which led to the development of COBOL,
            an early high-level programming language still in use today.`,
                bornAt: '1815-12-10T00:00:00.000Z',
                diedAt: '1852-11-27T00:00:00.000Z'
              }
            ]
          }
        },
        {
          date: "21:33 04.09.2022",
          data: {
            test: "test"
          }
        },
        {
          date: "21:33 04.09.2022",
          data: {
            test: "test"
          }
        }
      ]
    }
  },
  methods: {
    getChannels: function () {
      axios.get('/api/v1/channels/').then(response => {
        this.menu = response.data
      })
    },
    connect: function () {
      let socket = new WebSocket("ws://" + document.location.host + "/ws/root");

      socket.onmessage = event => {
        const response = JSON.parse(event.data)
        switch (response.event) {
          case 'UpdateChannels':
            this.getChannels()
            break;
          case 'UpdateCounts':
            this.getChannels()
        }

        console.log('Server:', response)
      };

      socket.onerror = event => {
        console.log("ERROR!!", event.data)
        socket.close()
      }
    },
  },
  mounted() {
    this.connect()
    this.getChannels()
  }
}
</script>

<style lang="scss" scoped>
  .wait {
    width: 100%;
    height: 100%;
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    font-size: 2em;
    text-transform: uppercase;
    .icon {
      font-size: 5em;
      color: #20c984;
      text-shadow: 0 0 20px #20c984;
    }
  }
</style>