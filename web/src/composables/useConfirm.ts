import { ref } from 'vue'

export interface ConfirmOptions {
  title: string
  message: string
  confirmText?: string
  cancelText?: string
  variant?: 'danger' | 'warning' | 'info'
}

const visible = ref(false)
const options = ref<ConfirmOptions>({
  title: '',
  message: '',
  confirmText: 'Confirm',
  cancelText: 'Cancel',
  variant: 'danger',
})

let resolvePromise: ((value: boolean) => void) | null = null

export function useConfirm() {
  function confirm(opts: ConfirmOptions): Promise<boolean> {
    options.value = {
      confirmText: 'Confirm',
      cancelText: 'Cancel',
      variant: 'danger',
      ...opts,
    }
    visible.value = true
    return new Promise<boolean>((resolve) => {
      resolvePromise = resolve
    })
  }

  function handleConfirm() {
    visible.value = false
    resolvePromise?.(true)
    resolvePromise = null
  }

  function handleCancel() {
    visible.value = false
    resolvePromise?.(false)
    resolvePromise = null
  }

  return {
    visible,
    options,
    confirm,
    handleConfirm,
    handleCancel,
  }
}
