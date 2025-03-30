// CommentModal.vue
<template>
    <div v-if="show" class="modal">
        <div class="modal-content">
          <span class="close" @click="$emit('close')">&times;</span>
          <slot name="header"></slot>
          <div class="modal-body">
               <div class="emoji-grid">
                  <div v-for="emoji in emojis" :key="emoji" class="emoji" @click="selectEmoji(emoji)">
                    {{ emoji }}
                  </div>
                </div>
              <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
          </div>
        </div>  
    </div>
  </template>

<script>
export default {
  props: {
    show: Boolean,
    msg: Object
  },
  data() {
    return {
      emojis: ['❤️', '🔥', '👍', '😂', '😭'],
      errormsg: '',
    };
  },
  methods: {
    selectEmoji(emoji) {
      this.$emit('emoji-selected', emoji);
    }
  }
};
</script>

<style scoped>
.modal {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.4);
  display: flex;
  justify-content: center;
  align-items: center;
}

.modal-content {
  background: #ffffff;
  border-radius: 20px;
  padding: 20px;
  width: 300px;
  box-shadow: 0px 4px 10px rgba(0, 0, 0, 0.1);
  text-align: center;
}

.close {
  float: right;
  font-size: 24px;
  cursor: pointer;
}

.emoji-grid {
  display: flex;
  justify-content: center;
  gap: 10px;
}

.emoji {
  font-size: 24px;
  cursor: pointer;
}
</style>
