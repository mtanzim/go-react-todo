import { useEffect, useState } from "react";
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

function addTodo(name: string): Promise<Todo> {
  return fetch(`${API_ENDPOINT}`, {
    method: "POST",
    body: JSON.stringify({ name }),
  }).then((res) => res.json());
}

function Todo() {
  const [todos, setTodos] = useState<Todo[]>([]);
  const [errMsg, setErrMsg] = useState<string | null>(null);
  const [loading, setLoading] = useState<boolean>(false);
  const [newTodo, setNewTodo] = useState<string>("");

  useEffect(() => {
    const onMount = async () => {
      try {
        setLoading(true);
        const todos = await fetchTodos();
        setTodos(todos);
      } catch (err) {
        if (err instanceof Error) {
          setErrMsg(err.message);
          return;
        }
        setErrMsg("Something went wrong");
      } finally {
        setLoading(false);
      }
    };
    onMount();
  }, []);

  async function handleAdd() {
    if (!newTodo) {
      return;
    }
    try {
      setLoading(true);
      const res = await addTodo(newTodo);
      setTodos((cur) => cur.concat(res));
    } catch (err) {
      if (err instanceof Error) {
        setErrMsg(err.message);
        return;
      }
      setErrMsg("Something went wrong");
    } finally {
      setLoading(false);
    }
  }

  if (loading) {
    return <div>Loading...</div>;
  }

  return (
    <div>
      <h1>To do list</h1>
      {todos.length > 0 ? (
        <ul>
          {todos.map((todo) => (
            <li>{todo.name}</li>
          ))}
        </ul>
      ) : (
        <div>No todos</div>
      )}
      <input type="text" onChange={(e) => setNewTodo(e.target.value)}></input>
      <button onClick={handleAdd}>Add</button>
      {errMsg ?? <p>{errMsg}</p>}
    </div>
  );
}

function App() {
  return <Todo />;
}

export default App;
