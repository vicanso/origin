<template>
  <header class="header">
    <div class="userInfo">
      <span v-if="fetchingUserInfo">正在加载...</span>
      <div class="functions" v-else-if="userAccount">
        <router-link :to="{ name: profile }" class="account">
          <i class="el-icon-user" />
          <span>{{ userAccount }}</span>
        </router-link>
        <span class="divided">|</span>
        <a class="logout" href="#" title="退出登录" @click="logout">
          <i class="el-icon-switch-button" />
        </a>
      </div>
      <div v-else>
        <router-link :to="{ name: login }" class="login">
          <i class="el-icon-user" />
          登录
        </router-link>
        <span class="divided">|</span>
        <router-link :to="{ name: register }" class="register">
          <i class="el-icon-circle-plus" />
          注册
        </router-link>
      </div>
    </div>
  </header>
</template>
<script>
import { mapState, mapActions } from "vuex";
import { LOGIN, REGISTER, PROFILE } from "@/constants/route";
export default {
  name: "MainHeader",
  computed: mapState({
    fetchingUserInfo: state => state.user.processing,
    userAccount: state => state.user.info.account
  }),
  data() {
    return {
      login: LOGIN,
      register: REGISTER,
      profile: PROFILE
    };
  },
  methods: {
    ...mapActions(["logout"])
  }
};
</script>
<style lang="sass" scoped>
@import "@/common.sass"
.header
  height: $mainHeaderHeight
  background-color: $white
  padding: 5px 0
  line-height: $mainHeaderHeight - 10
  color: $darkBlue
  box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08)
.userInfo
  float: right
  font-size: 13px
  margin-right: $mainMargin
  i
    margin-right: 3px
    font-weight: bold
.divided
  margin: 0 15px
</style>
