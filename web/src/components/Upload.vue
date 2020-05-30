<template>
  <div class="upload">
    <el-upload
      list-type="picture"
      drag
      :action="action"
      :on-success="handleSuccess"
      :on-error="handleError"
      :on-remove="handleRemove"
      :before-upload="handleBeforeUpload"
      :on-exceed="handleExceed"
      :limit="$props.limit"
      :file-list="fileList"
    >
      <i class="el-icon-upload"></i>
      <div class="el-upload__text">将文件拖到此处，或<em>点击上传</em></div>
    </el-upload>
  </div>
</template>
<script>
import { contains } from "@/helpers/util";
import { FILES } from "@/constants/url";
export default {
  props: {
    bucket: {
      type: String,
      required: true
    },
    limit: {
      type: Number,
      default: 1
    },
    files: Array
  },
  data() {
    const { files } = this.$props;
    let fileList = [];
    if (files) {
      fileList = files.slice(0);
    }
    return {
      fileList,
      action: FILES + "?bucket=" + this.$props.bucket
    };
  },
  methods: {
    handleRemove(file) {
      const fileList = this.fileList.filter(item => {
        if (file.url && item.url === file.url) {
          return false;
        }
        return item.url !== file.response.url;
      });
      this.fileList = fileList;
      this.$emit("change", this.fileList.slice(0));
    },
    handleSuccess(file) {
      this.fileList.push(file);
      this.$emit("change", this.fileList.slice(0));
    },
    handleError(err) {
      this.$message.error(err.message);
    },
    handleBeforeUpload(file) {
      if (!contains(["image/jpeg", "image/png"], file.type)) {
        this.$message.warning("仅支持上传JPG与PNG格式");
        return false;
      }
      const tooLarge = file.size / 1024 / 1024 > 1;

      if (tooLarge) {
        this.$message.error("上传图片不能超过1MB");
        return false;
      }
      return true;
    },
    handleExceed() {
      const { limit } = this.$props;
      this.$message.warning(`图片限制上传${limit}张，请先删除图片再上传`);
    }
  }
};
</script>
