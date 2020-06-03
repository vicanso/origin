import request from "@/helpers/request";

import { SUPPLIERS, SUPPLIERS_ID } from "@/constants/url";
import { attachUpdatedAtDesc, attachStatusDesc } from "@/store/modules/common";
import { addNoCacheQueryParam, findByID } from "@/helpers/util";

const prefix = "supplier";
const mutationSupplierProcessing = `${prefix}.processing`;
const mutationSupplierList = `${prefix}.list`;
const mutationSupplierUpdate = `${prefix}.update`;

const state = {
  processing: false,
  list: {
    data: null,
    count: -1
  }
};

export default {
  state,
  mutations: {
    [mutationSupplierProcessing](state, processing) {
      state.processing = processing;
    },
    [mutationSupplierList](state, { suppliers = [], count }) {
      if (count >= 0) {
        state.list.count = count;
      }
      suppliers.forEach(attachUpdatedAtDesc);
      state.list.data = suppliers;
    },
    [mutationSupplierUpdate](state, { id, data }) {
      if (!state.list.data) {
        return;
      }
      const found = findByID(state.list.data, id);
      if (found) {
        Object.assign(found, data);
        attachUpdatedAtDesc(found);
        attachStatusDesc(found);
      }
    }
  },
  actions: {
    async addSupplier({ commit }, supplier) {
      commit(mutationSupplierProcessing, true);
      try {
        const { data } = await request.post(SUPPLIERS, supplier);
        return data;
      } finally {
        commit(mutationSupplierProcessing, false);
      }
    },
    // listSupplier 获取供应商
    async listSupplier({ commit }, params) {
      commit(mutationSupplierProcessing, true);
      try {
        const { data } = await request.get(SUPPLIERS, {
          params: addNoCacheQueryParam(params)
        });
        commit(mutationSupplierList, data);
        return data;
      } finally {
        commit(mutationSupplierProcessing, false);
      }
    },
    // getSupplierByID 通过id获取supplier信息
    async getSupplierByID({ commit }, id) {
      const found = findByID(state.list.data, id);
      if (found) {
        return found;
      }
      commit(mutationSupplierProcessing, true);
      try {
        const url = SUPPLIERS_ID.replace(":id", id);
        const { data } = await request.get(url, {
          params: addNoCacheQueryParam()
        });
        attachUpdatedAtDesc(data);
        return data;
      } finally {
        commit(mutationSupplierProcessing, false);
      }
    },
    // updateSupplierByID 通过ID更新supplier信息
    async updateSupplierByID({ commit }, { id, data }) {
      if (!data || Object.keys(data).length === 0) {
        return;
      }
      commit(mutationSupplierProcessing, true);
      try {
        const url = SUPPLIERS_ID.replace(":id", id);
        await request.patch(url, data);
        commit(mutationSupplierUpdate, {
          id,
          data
        });
      } finally {
        commit(mutationSupplierProcessing, false);
      }
    }
  }
};
