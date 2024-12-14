<script lang="ts">
   import { ModeWatcher } from "mode-watcher";
    import { Button } from "$lib/components/ui/button/index.js";
    import { Input } from "$lib/components/ui/input/index.js";
    import { Label } from "$lib/components/ui/label/index.js";
    import * as Select from "$lib/components/ui/select";
  </script>
  
  <main>
<div class="flex min-h-screen w-full flex-col">
  

  <div class="w-full lg:grid lg:min-h-[600px] lg:grid-cols-2 xl:min-h-[800px]">
    <div class="flex items-center justify-center py-12">
      <div class="mx-auto grid w-[350px] gap-6">
        <div class="grid gap-2 text-center">
          <h1 class="text-3xl font-bold">Sign Up</h1>
          <p class="text-muted-foreground text-balance">
            Create your account today.
          </p>
        </div>
        <div class="grid gap-4">
          <div class="grid gap-2">
            <Label for="email">Phone</Label>
            <Input id="email" type="phone" placeholder="+91 99952 83835" required />
          </div>
          <div class="grid gap-2">
            <Label for="phone">Farmer ID</Label>
            <Input id="phone" type="phone" placeholder="777" required />
          </div>
          <div class="grid gap-2">
            <div class="flex items-center">
              <Label for="password">Password</Label>
            </div>
            <Input id="password" type="password" required />
          </div>
          <div class="grid gap-2">
            <div class="flex items-center">
              <Label for="select">User</Label>
            </div>
            <Select.Root>
              <Select.Trigger class="w-full">
                <Select.Value placeholder="User Type" />
              </Select.Trigger>
              <Select.Content>
                <Select.Item value="Farmer">Farmer</Select.Item>
                <Select.Item value="Buyer">System</Select.Item>
              </Select.Content>
            </Select.Root>
          </div>
          <Button type="submit" class="w-full">Login</Button>
        </div>
        <div class="mt-4 text-center text-sm">
          Already have an account?
          <a href="/login" class="underline"> Sign in </a>
        </div>
      </div>
    </div>
    <div class="bg-muted hidden lg:block">
      <img
        src="https://www.openaccessgovernment.org/wp-content/uploads/2021/02/dreamstime_xxl_183417536.jpg"
        alt="placeholder"
        width="1920"
        height="1080"
        class="h-full w-full object-cover dark:brightness-[0.2] dark:grayscale"
      />
    </div>
  </div>

  </div>
  </main>
