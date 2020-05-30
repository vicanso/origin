<template>
  <!-- 地区信息修改 -->
  <el-card>
    <div slot="header">
      <i class="el-icon-position" />
      更新地区信息
    </div>
    <!-- <el-form :model="currentRegion" label-width="80px" v-if="currentRegion">
      <el-row :gutter="15">
        <el-col :span="8">
          <el-form-item label="名称：">
            <el-input v-model="currentRegion.name" />
          </el-form-item>
        </el-col>
        <el-col :span="8">
          <el-form-item label="代码：">
            <el-input v-model="currentRegion.code" disabled />
          </el-form-item>
        </el-col>
        <el-col :span="8">
          <el-select
            class="selector"
            v-model="currentRegion.status"
            placeholder="请选择状态"
          >
            <el-option
              v-for="item in statuses"
              :key="item.value"
              :label="item.name"
              :value="item.value"
            />
          </el-select>
        </el-col>
        <el-col :span="12">
          <el-form-item>
            <el-button class="submit" type="primary" @click="update"
              >更新</el-button
            >
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item>
            <el-button class="submit" @click="goBack">返回</el-button>
          </el-form-item>
        </el-col>
      </el-row>
    </el-form> -->
  </el-card>
</template>
