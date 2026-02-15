<script setup lang="ts">
import { zodResolver } from '@primevue/forms/resolvers/zod';
import { z } from 'zod';
definePageMeta({
  layout: "guest",
  middleware: "guest",
});

const toast = useToast()
const auth = useAuthStore()
const router = useRouter()

const checked = ref(false)
const loading = ref(false)

const initialValues = reactive({
  email: "",
    Password: "",
});

const resolver = ref(zodResolver(
    z.object({
        email: z.email({ message: 'Invalid email address.' }),
        password: z.string().min(1, { message: 'Password is required.' }).min(6, { message: 'Password must be at least 8 characters.' })

    })
));

const handleSubmit = async ({ valid, values }:any) => {
  if (!valid) return;
  loading.value = true;
  try {
      console.log(values, checked.value);
    await new Promise((resolve) => setTimeout(resolve, 1000));

    toast.add({
      severity: 'success',
      summary: 'Link Sent',
      detail: `Instructions sent to ${values.email}`,
      life: 4000,
    });
    router.push('/')
  } finally {
    loading.value = false;
  }
};
</script>
<template>
    <div class="card w-full max-w-md p-6 sm:p-8 rounded-2xl shadow-sm hover:shadow-md flex flex-col gap-4 sm:gap-6 bg-surface-0 dark:bg-surface-900">

      <div class="flex flex-col items-center gap-2 text-center">
        <NuxtLink to="/">
          <img src="/favicon.ico" class="h-16 w-16 sm:h-20 sm:w-20" alt="Logo" />
        </NuxtLink>
        <h1 class="text-2xl sm:text-3xl font-semibold text-surface-900 dark:text-surface-0">Welcome Back</h1>
        <p class="text-muted-color text-sm">Sign in to continue</p>
      </div>

      <Form
        v-slot="$form"
        :initialValues
        :resolver
        @submit="handleSubmit"
        class="flex flex-col gap-4 sm:gap-6 w-full"
      >
        <div class="flex flex-col gap-1 w-full">
          <label for="email" class="text-sm font-medium text-surface-700 dark:text-surface-200">Email</label>
          <InputText id="email" name="email" type="text" placeholder="your@email.com" fluid />
          <Message v-if="$form.email?.invalid" severity="error" size="small" variant="simple">
            {{ $form.email.error?.message }}
          </Message>
        </div>

        <div class="flex flex-col gap-1 w-full">
          <label for="password" class="text-sm font-medium text-surface-700 dark:text-surface-200">Password</label>
          <Password
            id="password"
            name="password"
            placeholder="Enter your password"
            :feedback="false"
            toggleMask
            fluid
          />
          <Message v-if="$form.password?.invalid" severity="error" size="small" variant="simple">
            {{ $form.password.error?.message }}
          </Message>
        </div>

        <div class="flex flex-col sm:flex-row items-baseline sm:items-center justify-between w-full gap-2 sm:gap-0 text-sm">
          <div class="flex items-center">
            <Checkbox v-model="checked" id="rememberme" binary class="mr-2" />
            <label for="rememberme">Remember me</label>
          </div>
        </div>

        <Button type="submit" severity="primary" label="Sign In" :loading="loading" class="w-full rounded-full py-3 text-lg" />
      </Form>
    </div>
</template>
