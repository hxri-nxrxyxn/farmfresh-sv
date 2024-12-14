<script>
    import {Route, Router} from "svelte-routing"
    import Dash from "./Dash.svelte"
    import Login from "./Login.svelte"
    import Signup from "./Signup.svelte"
    import Orders from "./Orders.svelte"
    import Products from "./Products.svelte"
</script>

<Router>
    <Route path="/dashboard" component={Dash}/>
    <Route path="/login" component={Login}/>
    <Route path="/signup" component={Signup}/>
    <Route path="/orders" component={Orders}/>
    <Route path="/products" component={Products}/>
</Router>


<style></style>