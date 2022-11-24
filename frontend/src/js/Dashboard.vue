<template>
  <div class="dashboard-container">
    <transition-group tag="div" name="slide-fade" class="menu">
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
    <div class="logs">
      <div v-if="!this.connection" class="wait">
        <div class="icon">
          <font-awesome-icon icon="fa-solid fa-water"/>
        </div>
        <div class="text">
          Select channel
        </div>
      </div>
      <transition-group tag="div" name="slide-fade" v-else>
        <log-item v-for="log in logs"
                  :key="log.id"
                  :date="log.created_at"
                  :request="log.request"
                  :headers="log.headers"
                  :query="log.query"
                  :path="log.path"
        ></log-item>
      </transition-group>
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
      logs: []
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
            console.log('ADD!')
            this.logs.unshift(response.data)
        }
      };
    }
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
}
</style>