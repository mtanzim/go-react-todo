import { useEffect, useState } from "react";
import "./App.css";

interface Todo {
  id: number;
  name: string;
  completed: boolean;
}

type UpdateTodo = Omit<Todo, "id">;

const API_ENDPOINT = `http://localhost:8080/api/v1/todo`;

function fetchTodos(): Promise<Todo[]> {
  return fetch(`${API_ENDPOINT}`).then((res) => res.json());
}

function addTodo(name: string): Promise<Todo> {
  return fetch(`${API_ENDPOINT}`, {
    method: "POST",
    body: JSON.stringify({ name }),
  }).then((res) => res.json());
}

function updateTodo(id: number, todo: UpdateTodo): Promise<Todo> {
  return fetch(`${API_ENDPOINT}/${id}`, {
    method: "PUT",
    body: JSON.stringify(todo),
  }).then((res) => res.json());
}

function deleteTodo(id: number): Promise<void> {
  return fetch(`${API_ENDPOINT}/${id}`, {
    method: "DELETE",
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
  async function handleToggle(id: number, todo: Todo) {
    try {
      setLoading(true);
      const updated = { ...todo, completed: !todo.completed, id: undefined };
      const res = await updateTodo(id, updated);
      setTodos((cur) => cur.map((todo) => (todo.id === id ? res : todo)));
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

  async function handleDelete(id: number) {
    try {
      setLoading(true);
      await deleteTodo(id);
      setTodos((cur) => cur.filter((todo) => todo.id !== id));
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
        <div
          style={{
            display: "grid",
            gridTemplate: "1fr 1fr 1fr",
            gap: "4px",
            margin: "12px",
          }}
        >
          {todos.map((todo) => (
            <div
              style={{ display: "flex", gap: "4px", margin: "4px" }}
              key={todo.id}
            >
              <p
                style={
                  todo.completed
                    ? { textDecoration: "line-through" }
                    : { textDecoration: "none" }
                }
              >
                {todo.name}
              </p>
              <button onClick={() => handleDelete(todo.id)}>Delete</button>
              <button onClick={() => handleToggle(todo.id, todo)}>
                {!todo.completed ? "Mark Done" : "Do again"}
              </button>
            </div>
          ))}
        </div>
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
