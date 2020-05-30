<template>
  <el-card class="product">
    <div slot="header">
      <i class="el-icon-files" />
      添加/更新产品信息
    </div>

    <div class="upload">
      <h5>产品图片</h5>
      <Upload
        :files="files"
        bucket="origin-pics"
        @change="handleUpload"
        v-if="!processing"
      />
    </div>

    <el-form label-width="90px" class="form " v-loading="processing">
      <el-row :gutter="15">
        <el-col :span="8">
          <el-form-item label="名称：">
            <el-input v-model="form.name" placeholder="请输入产品名称" />
          </el-form-item>
        </el-col>
        <el-col :span="8">
          <el-form-item label="单价：">
            <el-input
              type="number"
              v-model="form.price"
              placeholder="请输入产品单价"
            />
          </el-form-item>
        </el-col>
        <el-col :span="8">
          <el-form-item label="单位：">
            <el-input v-model="form.unit" placeholder="请输入产品单位" />
          </el-form-item>
        </el-col>
        <el-col :span="8">
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
        <el-col :span="8">
          <el-form-item label="分类：">
            <el-select
              class="selector"
              multiple
              v-model="form.categories"
              placeholder="请选择分类"
            >
              <!-- TODO 产品分类拉取  -->
              <el-option key="a" label="a" value="a" />
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :span="8">
          <el-form-item label="品牌：">
            <BrandSelect @change="handleChangeBrand" />
          </el-form-item>
        </el-col>
        <el-col :span="8">
          <el-form-item label="产地：">
            <RegionSelect :maxLevel="2" :showAllLevels="true" />
          </el-form-item>
        </el-col>
        <el-col :span="8">
          <el-form-item label="SN：">
            <el-input v-model="form.sn" placeholder="请输入产品SN" />
          </el-form-item>
        </el-col>
        <el-col :span="8">
          <el-form-item label="主图："> </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="销售时间：">
            <el-date-picker
              class="dateRange"
              v-model="dateRange"
              type="daterange"
              range-separator="至"
              start-placeholder="开始日期"
              end-placeholder="结束日期"
            >
            </el-date-picker>
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="关键字：">
            <el-input
              v-model="form.keywords"
              placeholder="请输入产品关键字，多个关键字以空格分开"
            />
          </el-form-item>
        </el-col>
        <el-col :span="24">
          <el-form-item label="简介：">
            <el-input
              v-model="form.catalog"
              type="textarea"
              :autosize="{ minRows: 5, maxRows: 10 }"
              placeholder="请输入产品简介"
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
import { mapState, mapActions } from "vuex";
import BrandSelect from "@/components/products/BrandSelect.vue";
import Upload from "@/components/Upload.vue";
import RegionSelect from "@/components/region/Select.vue";

export default {
  name: "Product",
  components: {
    RegionSelect,
    BrandSelect,
    Upload
  },
  data() {
    return {
      submitText: "添加",
      processing: false,
      dateRange: null,
      files: null,
      form: {
        name: "",
        price: null,
        unit: "",
        catalog: "",
        pics: null,
        mainPic: null,
        sn: "",
        status: null,
        keywords: "",
        categories: null,
        startedAt: null,
        endedAt: null,
        origin: "",
        brand: null
      }
    };
  },
  computed: {
    ...mapState({
      statuses: state => state.product.statuses || [],
      listStatusProcessing: state => state.product.statusesListProcessing
    })
  },
  methods: {
    ...mapActions(["listProductStatus", "listBrand"]),
    handleChangeBrand(brand) {
      this.form.brand = brand;
    },
    handleUpload(files) {
      // TODO 记录至pics中
      console.dir(files);
    },
    async submit() {},
    goBack() {
      this.$router.back();
    }
  },
  async beforeMount() {
    try {
      await this.listProductStatus();
    } catch (err) {
      this.$message.error(err.message);
    }
  }
};
</script>
<style lang="sass" scoped>
$uploadWidth: 350px
.product
  .upload
    float: right
    width: $uploadWidth
  h5
    margin: 0 0 10px 0
  .form
    margin-right: $uploadWidth + 10
  .selector, .btn, .dateRange
    width: 100%
</style>
