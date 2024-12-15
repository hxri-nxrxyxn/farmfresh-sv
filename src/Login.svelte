<script lang="ts">
    import { Button } from "$lib/components/ui/button/index.js";
    import { Input } from "$lib/components/ui/input/index.js";
    import { Label } from "$lib/components/ui/label/index.js";

  let username = $state('');
  let password = $state('');
  let isLoggedIn = false;

  let data = {
  "username": username,
  "password": password
};




  async function login() {
    let jsonData = JSON.stringify(data);
    

    try {
      const response = await fetch('http://127.0.0.1:8080/v1/users/login/', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },

        body: JSON.stringify({
          username,
          password
        })
      })

      console.log("Data ",jsonData);
      console.log(response);
      isLoggedIn = true;
      errorMessage = '';

      if (response.ok) {
        const data = await response.json();
        // Handle successful login, e.g., set authentication token, redirect
        isLoggedIn = true;
      } else {
<<<<<<< HEAD
        console.log("hai 1");
        //const errorData = await response.json();
        console.log("hai 2");
        //errorMessage = errorData.message || 'Login failed';
        errorMessage = 'Login failed';
=======
        const errorData = await response.json();
>>>>>>> 50bced5f8bddd412c54f02d4941135f4f965dc4d
      }
    } catch (error) {
      console.log("SHHHHH");
      console.error('Error logging in:', error);
    }
  }
<<<<<<< HEAD
  //login()
=======


>>>>>>> 50bced5f8bddd412c54f02d4941135f4f965dc4d
  </script>

  <button id="s" onclick={login}>heyy</button>
  
  <div class="w-full lg:grid lg:min-h-[600px] lg:grid-cols-2 xl:min-h-[800px]">
    <div class="flex items-center justify-center py-12">
      <div class="mx-auto grid w-[350px] gap-6">
        <div class="grid gap-2 text-center">
          <h1 class="text-3xl font-bold">Login</h1>
          <p class="text-muted-foreground text-balance">
            Continue with your phone number
          </p>
        </div>
        <div class="grid gap-4">
          <div class="grid gap-2">
            <Label for="phone">Phone</Label>
            <Input id="email" type="name" bind:value={username} placeholder="+91 99952 83835" required />
          </div>
          <div class="grid gap-2">
            <div class="flex items-center">
              <Label for="password" id="getpass">Password</Label>
            </div>
            <Input id="password" type="password" required />
          </div>
          <a href="/dashboard">
          <Button {onclick} type="submit" class="w-full">Login</Button>
        </a>
        </div>
        <div class="mt-4 text-center text-sm">
          Don&apos;t have an account?

          <a href="/signup" class="underline"> Sign up </a>
        </div>
      </div>
    </div>
    <div class="bg-muted hidden lg:block">
      <img
        src="https://www.openaccessgovernment.org/wp-content/uploads/2021/02/dreamstime_xxl_183417536.jpg"
        alt="background"
        width="1920"
        height="1080"
        class="h-full w-full object-cover dark:brightness-[0.2] dark:grayscale"
      />
    </div>
  </div>
