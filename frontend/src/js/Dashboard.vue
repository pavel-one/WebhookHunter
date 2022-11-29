<template>
  <div class="description">
    Send any request <br>
    <span class="bg">
      {{ currentUrl }}/{any*}
    </span>
  </div>
  <div class="dashboard-container">
    <transition-group tag="div" name="slide-fade" class="menu" :class="{'selected': this.connection}">
      <menu-item v-for="item in menu"
                 :key="item.id"
                 :id="item.id"
                 :name="item.path"
                 :count="item.count"
                 :active="item.active"
                 @changeChannel="changeChannel"
                 :date="item.date">
      </menu-item>
    </transition-group>
    <div class="logs" :class="{'selected': this.connection}">
      <div v-if="!this.connection" class="wait">
        <div class="icon">
          <svg xmlns="http://www.w3.org/2000/svg" class="icon icon-tabler icon-tabler-ripple" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" fill="none" stroke-linecap="round" stroke-linejoin="round">
            <path stroke="none" d="M0 0h24v24H0z" fill="none"></path>
            <path d="M3 7c3 -2 6 -2 9 0s6 2 9 0"></path>
            <path d="M3 17c3 -2 6 -2 9 0s6 2 9 0"></path>
            <path d="M3 12c3 -2 6 -2 9 0s6 2 9 0"></path>
          </svg>
        </div>
        <div class="text">
          SELECT CHANNEL
        </div>
      </div>
      <div v-else-if="!this.logs" class="wait">
        <div class="icon">
          <svg xmlns="http://www.w3.org/2000/svg" class="icon icon-tabler icon-tabler-player-play" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" fill="none" stroke-linecap="round" stroke-linejoin="round">
            <path stroke="none" d="M0 0h24v24H0z" fill="none"></path>
            <path d="M7 4v16l13 -8z"></path>
          </svg>
        </div>
        <div class="text">
          Send any request to <br>
          <span class="bg">{{ currentUrl }}</span>
        </div>
      </div>
      <div v-else>
        <log-item v-for="log in logs"
                  :key="log.id"
                  :date="log.created_at"
                  :request="log.request"
                  :headers="log.headers"
                  :query="log.query"
                  :path="log.path"
        ></log-item>
      </div>
    </div>
  </div>
</template>

<script>
import MenuItem from "./chunks/menu-item";
import LogItem from "./chunks/log-item";
import axios from "axios";
import Logo from "./components/Logo";

export default {
  components: {Logo, MenuItem, LogItem},
  data: function () {
    return {
      connection: null,
      menu: [],
      logs: [],
      currentUrl: ''
    }
  },
  methods: {
    getChannels: function () {
      axios.get('/api/v1/channels/').then(response => {
        if (!response.data) {
          return
        }
        this.menu = response.data
      })
    },
    updateCount: function (data) {
      this.menu.forEach((item, index) => {
        if (item.id === data.id) {
          this.menu[index].count = data.count
        }
      })
    },
    addChannel: function (data) {
      if (!this.menu) {
        this.menu = []
      }

      this.menu.unshift(data)
    },
    dropChannel: function (data) {
      this.menu.forEach((item, index) => {
        if (item.id === data.id) {
          console.log('dropped', item.id, data.id)
          this.menu.splice(index, 1)
        }
      })
    },
    connect: function () {
      let wsProtocol = "ws://"
      if (window.location.protocol === 'https:') {
        wsProtocol = 'wss://'
      }

      let socket = new WebSocket(wsProtocol + document.location.hostname + ":8080/root");

      socket.onmessage = event => {
        const response = JSON.parse(event.data)
        switch (response.event) {
          case 'AddChannel':
            this.addChannel(response.data)
            break
          case 'UpdateCount':
            this.updateCount(response.data)
            break
          case 'DropChannel':
            this.dropChannel(response.data)
            break
        }

        console.log('Server:', response)
      };

      socket.onerror = event => {
        console.log("ERROR!!", event.data)
        socket.close()
      }
    },
    changeChannel: function (channelName) {
      this.menu.forEach((item, index) => {
        if (item.path === channelName) {
          console.log(index, channelName, item)
          this.menu[index].active = true
        } else {
          this.menu[index].active = false
        }
      })

      let wsProtocol = "ws://"
      if (window.location.protocol === 'https:') {
        wsProtocol = 'wss://'
      }

      let url = wsProtocol + document.location.hostname + ":8080" + channelName
      if (channelName === '/') {
        url = wsProtocol + document.location.hostname + ":8080/"
      }

      if (this.connection) {
        console.log('close connection old')
        this.connection.close()
      }

      this.connection = new WebSocket(url);

      this.connection.onerror = event => {
        console.log("ERROR CHANNEL!!", event.data)
        this.connection.close()
      }

      this.connection.onmessage = event => {
        const response = JSON.parse(event.data)
        console.log('Server channel:', response)

        switch (response.event) {
          case 'Load':
            this.logs = response.data
            break
          case 'Add':
            if (!this.logs) {
              this.logs = []
            }

            this.logs.unshift(response.data)
        }
      };
    }
  },
  mounted() {
    this.connect()
    this.getChannels()
    this.currentUrl = window.location.href.replace('/ui/', '')
  }
}
</script>

<style lang="scss" scoped>
.bg {
  background: #000000;
  padding: 10px;
  border-radius: 5px;
  box-shadow: 0 2px 5px rgb(0 0 0 / 50%);
}
.description {
  font-size: 1.5em;
  margin: 10px;
  line-height: 1.5;
  text-align: center;
}
.wait {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  font-size: 2em;
  text-transform: uppercase;
  .text {
    text-align: center;
    line-height: 1.2;
    text-transform: none;
  }
}

.dashboard-logo {
  position: absolute;
  bottom: 0;
  width: 100%;
  left: 0;
  opacity: .5;
}

.icon {
  font-size: 5em;
  color: #20c984;
  text-shadow: 0 0 20px #20c984;

  .icon-tabler {
    width: 150px;
    height: 150px;
  }
}
</style>