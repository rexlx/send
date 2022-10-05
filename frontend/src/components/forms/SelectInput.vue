<template>
    <div class="mb-3">
        <label :for="name" class="form-label">{{label}}</label>
        <select :id="name"
        class="form-select"
        :name="name"
        :required="required"
        :multiple="multi"
        :value="mVal"
        @input="$emit('update:mVal', $event.target.value)">
            <option disabled value="">{{ first }}</option>
            <option v-for="option in items" :value="option.value" :key="option.value">
                {{option.text}}
            </option>
        </select>
    </div>
</template>

<script>
export default {
    name: 'SelectInput',
    props: ["items", "name", "required", "label", "mVal", "multi", "first"],
}
</script>

<style scoped>
/* select {
    background: rgb(30, 30, 40);
    border: black;
} */
</style>