<template>
  <el-card>
    <div slot="header">
      <i class="el-icon-goods" />
      添加/更新品牌信息
    </div>

    <div class="upload">
      <h5>品牌LOGO</h5>
      <Upload
        :files="files"
        bucket="origin-pics"
        @change="handleUpload"
        v-if="!processing"
      />
    </div>
    <el-form label-width="80px" class="form" v-loading="processing">
      <el-row :gutter="15">
        <el-col :span="12">
          <el-form-item label="名称：">
            <el-input v-model="form.name" placeholder="请输入产品名称" />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="状态：">
            <el-select
              class="selector"
              v-model="form.status"
              v-loading="listStatusProcessing"
            >
              <el-option
                v-for="item in statuses"
                :key="item.value"
                :label="item.name"
                :value="item.value"
              ></el-option>
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :span="24">
          <el-form-item label="简介：">
            <el-input
              v-model="form.catalog"
              type="textarea"
              :autosize="{ minRows: 5, maxRows: 10 }"
              placeholder="请输入品牌简介"
            />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item>
            <el-button class="btn" type="primary" @click="submit">{{
              submitText
            }}</el-button>
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item>
            <el-button class="btn" @click="goBack">返回</el-button>
          </el-form-item>
        </el-col>
      </el-row>
    </el-form>
  </el-card>
</template>
<script>
import { mapActions, mapState } from "vuex";
import Upload from "@/components/Upload.vue";
import { diff } from "@/helpers/util";

export default {
  name: "Brand",
  components: {
    Upload
  },
  data() {
    const pageSizes = [10, 20, 30, 50];
    return {
      pageSizes,
      processing: false,
      query: {
        offset: 0,
        limit: pageSizes[0],
        order: "-updatedAt"
      },
      submitText: "添加",
      files: null,
      originalData: null,
      form: {
        name: "",
        status: null,
        logo: "",
        catalog: ""
      }
    };
  },
  computed: {
    ...mapState({
      listStatusProcessing: state => state.brand.statusesListProcessing,
      statuses: state => state.brand.statuses
    }),
    currentPage() {
      const { offset, limit } = this.query;
      return Math.floor(offset / limit) + 1;
    }
  },
  methods: {
    ...mapActions([
      "listBrandStatus",
      "addBrand",
      "getBrandByID",
      "updateBrandByID"
    ]),
    handleUpload(files) {
      this.form.logo = files[0].url || "";
    },
    async submit() {
      const { originalData } = this;
      const { name, status, logo, catalog } = this.form;
      if (!name || !status || !logo || !catalog) {
        this.$message.warning("名称、状态、LOGO与简介均不能为空");
        return;
      }
      try {
        const data = {
          name,
          status,
          logo,
          catalog
        };
        if (originalData) {
          const update = diff(data, this.originalData);
          if (!update.modifiedCount) {
            this.$message.warning("未修改信息无需要更新");
            return;
          }
          await this.updateBrandByID({
            id: originalData.id,
            data: update.data
          });
          this.$message.info("已成功更新该品牌信息");
        } else {
          await this.addBrand(data);
          this.$message.info("已成功添加该品牌");
        }

        this.$router.back();
      } catch (err) {
        this.$message.error(err.message);
      }
    },
    goBack() {
      this.$router.back();
    }
  },
  async beforeMount() {
    const { id } = this.$route.query;
    this.processing = true;
    try {
      if (id) {
        const data = await this.getBrandByID(Number(id));
        this.submitText = "更新";
        this.originalData = data;
        Object.assign(this.form, data);
        this.files = [
          {
            url: data.logo
          }
        ];
      }
      await this.listBrandStatus();
    } catch (err) {
      this.$message.error(err.message);
    } finally {
      this.processing = false;
    }
  }
};
</script>
<style lang="sass" scoped>
$uploadWidth: 350px
.upload
  float: right
  width: $uploadWidth
.form
  margin-right: $uploadWidth + 10
h5
  margin: 0 0 10px 0
.selector, .btn
  width: 100%
</style>
