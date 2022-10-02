<template>
    <form 
        @submit.prevent="submit"
        :ref="name"
        :event="event"
        autocomplete="off" 
        :method="method" 
        :action="action" 
        class="needs-validation" 
        novalidate>
        <slot></slot>
    </form>
</template>

<script>
export default {
    name: 'FormTag',
    props: ["method", "action", "name", "event"],
    methods: {
        submit() {
            let myForm = this.$refs[this.$props.name];

            if (myForm.checkValidity()) {
                this.$emit(this.$props['event']);
            }
            myForm.classList.add("was-validated");
        }
    },
}
</script>