import { useEffect, useState } from "react";
import reactLogo from "./assets/react.svg";
import viteLogo from "/vite.svg";
import "./App.css";

interface Todo {
  id: number;
  name: string;
  completed: boolean;
}

const API_ENDPOINT = `http://localhost:8080/api/v1/todo/`;

function fetchTodos(): Promise<Todo[]> {
  return fetch(`${API_ENDPOINT}`).then((res) => res.json());
}

function Todo() {
  const [todos, setTodos] = useState<Todo[]>([]);
  const [errMsg, setErrMsg] = useState<string | null>(null);
  const [loading, setLoading] = useState<boolean>(false);
  useEffect(() => {
    const onMount = async () => {
      try {
        const todos = await fetchTodos();
        setTodos(todos);
      } catch (err) {
        if (err instanceof Error) {
          setErrMsg(err.message);
          return;
        }
        setErrMsg("Something went wrong");
      }
    };
    onMount();
  }, []);

  if (loading) {
    return <div>Loading...</div>;
  }

  if (errMsg) {
    return <div>{errMsg}</div>;
  }

  return (
    <div>
      <ul>
        {todos.map((todo) => (
          <li>{todo.name}</li>
        ))}
      </ul>
    </div>
  );
}

function App() {

  return (<Todo />);
}

export default App;
