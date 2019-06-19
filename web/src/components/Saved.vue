<template>
  <el-table
    :data="saved"
    :default-sort="{prop: 'status', order: 'descending'}"
    style="width: 100%"
    row-key="UUID"
    :expand-row-keys="expands"
    @row-click="rowClick"
  >
    <el-table-column type="expand">
      <template slot-scope="props">
        <el-button type="danger" icon="el-icon-delete" round @click="Remove(props.row.UUID)">Remove</el-button>
        <el-form label-position="left" label-width="80px" class="demo-table-expand">
          <el-form-item label="Memory">
            <span>{{ props.row.memoryUsed }}GiB / {{ props.row.memoryTotal }}GiB</span>
          </el-form-item>
          <el-form-item label="Swap">
            <span>{{ props.row.swapUsed }}GiB / {{ props.row.swapTotal }}GiB</span>
          </el-form-item>
          <el-form-item label="Disk">
            <span>{{ props.row.diskUsed }}GiB / {{ props.row.diskTotal }}GiB</span>
          </el-form-item>
        </el-form>
      </template>
    </el-table-column>
    <el-table-column label="Name" prop="serverName"></el-table-column>
    <el-table-column label="Status" prop="status" sortable>
      <template slot-scope="scope">
        <el-tag
          :type="scope.row.status === 'Online'? 'success' : 'danger'"
          disable-transitions
        >{{ scope.row.status }}</el-tag>
      </template>
    </el-table-column>
    <el-table-column label="IPv4" prop="IPv4">
      <template slot-scope="scope">
        <el-tag
          :type="scope.row.IPv4 === 'On'? 'success' : 'danger'"
          disable-transitions
        >{{ scope.row.IPv4 }}</el-tag>
      </template>
    </el-table-column>
    <el-table-column label="IPv6" prop="IPv6">
      <template slot-scope="scope">
        <el-tag
          :type="scope.row.IPv6 === 'On'? 'success' : 'danger'"
          disable-transitions
        >{{ scope.row.IPv6 }}</el-tag>
      </template>
    </el-table-column>
    <el-table-column label="System" prop="system">
      <template slot-scope="scope">
        <pre>{{ scope.row.system }}</pre>
      </template>
    </el-table-column>
    <el-table-column label="Location" prop="location"></el-table-column>
    <el-table-column label="Uptime" prop="uptime">
      <template slot-scope="scope">
        <pre>{{ scope.row.uptime }}</pre>
      </template>
    </el-table-column>
    <el-table-column label="Network" prop="network">
      <template slot-scope="scope">
        <p>⬆️{{ scope.row.uploadSpeed }}</p>
        <p>⬇️{{ scope.row.downloadSpeed }}</p>
      </template>
    </el-table-column>
    <el-table-column label="CPU" prop="CPUUsed">
      <template slot-scope="scope">
        <el-progress
          :text-inside="true"
          :stroke-width="26"
          :percentage="scope.row.CPUUsed"
          :color="customColors"
        ></el-progress>
      </template>
    </el-table-column>
    <el-table-column label="RAM" prop="memoryRate">
      <template slot-scope="scope">
        <el-progress
          :text-inside="true"
          :stroke-width="26"
          :percentage="scope.row.memoryRate"
          :color="customColors"
        ></el-progress>
      </template>
    </el-table-column>
    <el-table-column label="Disk" prop="diskRate">
      <template slot-scope="scope">
        <el-progress
          :text-inside="true"
          :stroke-width="26"
          :percentage="scope.row.diskRate"
          :color="customColors"
        ></el-progress>
      </template>
    </el-table-column>
  </el-table>
</template>

<script>
export default {
  data () {
    return {
      saved: null,
      customColors: [
        { color: '#32cd32', percentage: '50' },
        { color: '#ffa500', percentage: '85' },
        { color: '#ff0000', percentage: '100' }
      ],
      expands: [],
      timer: ''
    }
  },
  created: function () {
    this.fetchData()
    this.timer = setInterval(this.fetchData, 1000)
  },
  watch: {},
  methods: {
    fetchData () {
      this.$axios.get('/api/servers').then(res => {
        this.saved = []
        for (let server of res.data.servers) {
          var pushed = server
          var date = new Date().getTime()
          var timeStamp = Math.floor(date / 1000)
          if (timeStamp - pushed.lastTime > 30) {
            pushed.status = 'Offline'
            pushed.IPv4 = 'Off'
            pushed.IPv6 = 'Off'
            pushed.system = '-'
            pushed.uptime = '-'
            pushed.downloadSpeed = '-'
            pushed.uploadSpeed = '-'
            pushed.CPUUsed = NaN
            pushed.memoryRate = NaN
            pushed.diskRate = NaN
            pushed.memoryUsed = pushed.memoryTotal = pushed.swapUsed = pushed.swapTotal = pushed.diskUsed = pushed.diskTotal = 0
          } else {
            pushed.status = 'Online'
            pushed.IPv4 = pushed.IPv4Addr === '' ? 'Off' : 'On'
            pushed.IPv6 = pushed.IPv6Addr === '' ? 'Off' : 'On'
            pushed.memoryRate = parseInt(
              (pushed.memoryUsed / pushed.memoryTotal) * 100
            )
            pushed.diskRate = parseInt(
              (pushed.diskUsed / pushed.diskTotal) * 100
            )
            pushed.memoryUsed = pushed.memoryUsed.toFixed(2)
            pushed.memoryTotal = pushed.memoryTotal.toFixed(2)
            pushed.swapUsed = pushed.swapUsed.toFixed(2)
            pushed.swapTotal = pushed.swapTotal.toFixed(2)
            pushed.diskUsed = pushed.diskUsed.toFixed(2)
            pushed.diskTotal = pushed.diskTotal.toFixed(2)
          }
          // pushed.network = '⬆️' + pushed.uploadSpeed + '\n⬇️' + pushed.downloadSpeed
          this.saved.push(pushed)
        }
      })
    },
    Remove: function (msg) {
      this.$axios
        .post('/api/servers', {
          from: 'saved',
          UUID: msg
        })
        .then(res => {
          if (res.data === 'success') {
            this.fetchData()
            this.$message({
              showClose: true,
              message: 'Remove success!',
              type: 'success'
            })
            this.expands.splice(this.expands.indexOf(msg), 1)
          } else if (res.data === 'fail') {
            this.$message({
              showClose: true,
              message: 'Remove fail!',
              type: 'error'
            })
          }
        })
    },
    rowClick (row, event, column) {
      if (this.expands.indexOf(row.UUID) < 0) {
        this.expands.push(row.UUID)
      } else {
        this.expands.splice(this.expands.indexOf(row.UUID), 1)
      }
    }
  }
}
</script>

<style>
.demo-table-expand {
  font-size: 0;
  text-align: center;
}
.demo-table-expand label {
  width: 90px;
  color: #99a9bf;
}
.demo-table-expand .el-form-item {
  margin-left: 0;
  margin-bottom: 0;
  width: 50%;
}
.el-table .cell {
  white-space: pre-line;
}
</style>
