import request from "@/helpers/request";

import { BRANDS, BRANDS_ID } from "@/constants/url";
import {
  listStatus,
  attachStatusDesc,
  attachUpdatedAtDesc
} from "@/store/modules/common";
import { addNoCacheQueryParam, toUploadFiles, findByID } from "@/helpers/util";

const prefix = "brand";
const mutationBrandProcessing = `${prefix}.processing`;
const mutationBrandList = `${prefix}.list`;
const mutationBrandUpdate = `${prefix}.update`;

const state = {
  processing: false,
  list: {
    data: null,
    count: -1
  }
};

function enhanceBrandInfo(item) {
  attachStatusDesc(item);
  attachUpdatedAtDesc(item);
  item.files = toUploadFiles(item.logo);
}

export default {
  state,
  mutations: {
    [mutationBrandProcessing](state, processing) {
      state.processing = processing;
    },
    [mutationBrandList](state, { brands = [], count }) {
      if (count >= 0) {
        state.list.count = count;
      }
      brands.forEach(enhanceBrandInfo);
      state.list.data = brands;
    },
    [mutationBrandUpdate](state, { id, data }) {
      if (!state.list.data) {
        return;
      }
      const found = findByID(state.list.data, id);
      if (found) {
        Object.assign(found, data);
        enhanceBrandInfo(found);
      }
    }
  },
  actions: {
    // addBrand 添加品牌信息
    async addBrand({ commit }, brand) {
      commit(mutationBrandProcessing, true);
      try {
        const { data } = await request.post(BRANDS, brand);
        return data;
      } finally {
        commit(mutationBrandProcessing, false);
      }
    },
    listBrandStatus: listStatus,
    // listBrand 获取品牌
    async listBrand({ commit }, params) {
      commit(mutationBrandProcessing, true);
      try {
        await listStatus({ commit });
        const { data } = await request.get(BRANDS, {
          params: addNoCacheQueryParam(params)
        });
        commit(mutationBrandList, data);
        return data;
      } finally {
        commit(mutationBrandProcessing, false);
      }
    },
    // getBrandByID 通过id获取brand信息
    async getBrandByID({ commit }, id) {
      const found = findByID(state.list.data, id);
      if (found) {
        return found;
      }
      commit(mutationBrandProcessing, true);
      try {
        const url = BRANDS_ID.replace(":id", id);
        const { data } = await request.get(url, {
          params: addNoCacheQueryParam()
        });
        enhanceBrandInfo(data);
        return data;
      } finally {
        commit(mutationBrandProcessing, false);
      }
    },
    // updateBrandByID 通过id更新brand信息
    async updateBrandByID({ commit }, { id, data }) {
      if (data.files) {
        data.logo = data.files[0].url;
        delete data.files;
      }
      if (!data || Object.keys(data).length === 0) {
        return;
      }
      commit(mutationBrandProcessing, true);
      try {
        const url = BRANDS_ID.replace(":id", id);
        await request.patch(url, data);
        commit(mutationBrandUpdate, {
          id,
          data
        });
      } finally {
        commit(mutationBrandProcessing, false);
      }
    }
  }
};
