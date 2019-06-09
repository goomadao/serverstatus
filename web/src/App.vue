<template>
  <div id="app">
    <header class="grid-content header-color">
      <el-row>
        <a class="brand" href="#" :xs="24" :md="4">Server Status</a>
      </el-row>
    </header>

    <section>
      <el-col :gutter="20">
        <el-row id="side-nav">
          <el-menu
            background-color="#545c64"
            text-color="#fff"
            active-text-color="#ffd04b"
            :default-active="active"
            mode="horizontal"
            theme="light"
            router="false"
            @select="handleSelect"
          >
            <el-submenu index="/servers">
              <template slot="title">Servers</template>
              <el-menu-item index="/">Saved</el-menu-item>
              <el-menu-item index="/temp">Temp</el-menu-item>
            </el-submenu>
            <el-menu-item index>Help</el-menu-item>
          </el-menu>
        </el-row>

        <el-row>
          <div id="content">
            <router-view></router-view>
          </div>
        </el-row>
      </el-col>
    </section>
    <!-- <footer>Server Status</footer> -->
  </div>
</template>

<script>
export default {
  data () {
    return {
      active: '',
      abc: 123
    }
  },
  created: function () {
    var url = document.location.toString()
    var arrUrl = url.split('//')
    var start = arrUrl[1].indexOf('/')
    var relUrl = arrUrl[1].substring(start)
    if (relUrl.indexOf('?') !== -1) {
      relUrl = relUrl.split('?')[0]
    }
    if (relUrl[relUrl.length - 1] === '/') {
      this.active = '/'
    } else if (relUrl[relUrl.length - 1] === 'p') {
      this.active = '/temp'
    }
  },
  methods: {
    handleSelect (key, path) {
      if (key === '') {
        window.open('http://github.com/goomadao/serverstatus')
      }
    }
  }
}
</script>

<style>
body {
  background-color: #fafafa;
  margin: 0px;
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Oxygen,
    Ubuntu, Cantarell, "Open Sans", "Helvetica Neue", sans-serif;
}
header {
  width: 100%;
  height: 60px;
}
.header-color {
  background: #58b7ff;
}
#content {
  margin-top: 20px;
  padding-right: 40px;
}
.brand {
  color: #fff;
  background-color: transparent;
  margin-left: 20px;
  float: left;
  line-height: 25px;
  font-size: 25px;
  padding: 15px 15px;
  height: 30px;
  text-decoration: none;
}
</style>
