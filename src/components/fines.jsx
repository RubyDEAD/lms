// src/components/fines.jsx
import React, { useState, useEffect } from "react";
import axios from "axios";
import { supabase } from "../supabaseClient";
import { useNavigate } from "react-router-dom";
import { createClient } from "graphql-ws";
import "bootstrap/dist/css/bootstrap.min.css";

const Fines = () => {
    const [fines, setFines] = useState([]);
    const [loading, setLoading] = useState(true);
    const [authLoading, setAuthLoading] = useState(true);
    const [error, setError] = useState(null);
    const [isAuthenticated, setIsAuthenticated] = useState(false);
    const navigate = useNavigate();

    const API_URL = "http://localhost:8000/query";
    //const WS_URL = "ws://localhost:8000/query";

    useEffect(() => {
        const checkAuth = async () => {
            try {
                setAuthLoading(true);
                const { data: { session }, error } = await supabase.auth.getSession();
                if (error || !session) throw error;

                setIsAuthenticated(true);
                await fetchFines(session);
                setupSubscription(session.access_token);
            } catch (err) {
                console.error("Auth error:", err);
                setError("Session expired. Please log in again.");
                navigate("/login");
            } finally {
                setAuthLoading(false);
            }
        };

        checkAuth();

        return () => {
            if (window.fineSocket) {
                window.fineSocket.close();
            }
        };
    }, []);

    const fetchFines = async (session) => {
        try {
            setLoading(true);
            const token = session.access_token;

            const response = await axios.post(
                API_URL,
                {
                    query: `
                        query {
                            listFines {
                                fine_id
                                patronId
                                bookId
                                daysLate
                                ratePerDay
                                amount
                                createdAt
                                violationRecordId
                            }
                        }
                    `,
                },
                {
                    headers: {
                        Authorization: `Bearer ${token}`,
                    },
                }
            );

            console.log(response.data);

            const userFines = response.data.data.listFines.filter(
                (fine) => fine.patronId === session.user.id
            );

            setFines(userFines);
        } catch (err) {
            console.error("Error fetching fines:", err);
            setError("Failed to fetch fines.");
        } finally {
            setLoading(false);
        }
    };

        const setupSubscription = (token) => {
        const client = createClient({
            url: "ws://localhost:8000/query", // direct to fine-service
            connectionParams: {
                headers: {
                    Authorization: `Bearer ${token}`,
                },
            },
        });

        const unsubscribe = client.subscribe(
            {
                query: `
                    subscription {
                        fineCreated {
                            fine_id
                            patronId
                            bookId
                            daysLate
                            ratePerDay
                            amount
                            createdAt
                            violationRecordId
                        }
                    }
                `,
            },
            {
                next: ({ data }) => {
                    setFines((prev) => [...prev, data.fineCreated]);
                },
                error: (err) => console.error("Subscription error", err),
                complete: () => console.log("Subscription completed"),
            }
        );

        return unsubscribe; // optionally store for cleanup
    };

    const deleteFine = async (fineId) => {
        try {
            const { data: { session } } = await supabase.auth.getSession();
            const token = session.access_token;

            await axios.post(
                API_URL,
                {
                    query: `
                        mutation DeleteFine($fine_id: ID!) {
                            deleteFine(fine_id: $fine_id)
                        }
                    `,
                    variables: { fine_id: fineId },
                },
                {
                    headers: {
                        Authorization: `Bearer ${token}`,
                    },
                }
            );

            setFines((prev) => prev.filter((f) => f.fine_id !== fineId));
        } catch (err) {
            console.error("Error deleting fine:", err);
            setError("Failed to delete fine.");
        }
    };

    const createFine = async (e) => {
        e.preventDefault();
        try {
            const { data: { session } } = await supabase.auth.getSession();
            const token = session.access_token;

            const input = {
                patronId: session.user.id,
                bookId: e.target.bookId.value,
                daysLate: parseInt(e.target.daysLate.value),
                ratePerDay: parseFloat(e.target.ratePerDay.value),
                amount: parseInt(e.target.daysLate.value) * parseFloat(e.target.ratePerDay.value),
                violationType: "Late_Return", 
            };


                  await axios.post(
          API_URL,
          {
              query: `
                  mutation CreateFine(
                      $patronId: ID!,
                      $bookId: ID!,
                      $ratePerDay: Float!,
                      $violationType: ViolationType!,
                      $daysLate: Int
                  ) {
                      createFine(
                          patronId: $patronId,
                          bookId: $bookId,
                          ratePerDay: $ratePerDay,
                          violationType: $violationType,
                          daysLate: $daysLate
                      ) {
                          fine_id
                          patronId
                          bookId
                          daysLate
                          ratePerDay
                          amount
                          createdAt
                          violationRecordId
                      }
                  }
              `,
              variables: {
                  patronId: input.patronId,
                  bookId: input.bookId,
                  ratePerDay: input.ratePerDay,
                  violationType: "Late_Return", // or prompt user input
                  daysLate: input.daysLate,
              },
          },
          {
              headers: {
                  Authorization: `Bearer ${token}`,
              },
          }
      );


            e.target.reset();
            setError(null);
        } catch (err) {
            console.error("Error creating fine:", err);
            setError("Failed to create fine.");
        }
    };

    if (authLoading) return <div className="container mt-5">Checking authentication...</div>;
    if (!isAuthenticated) return <div className="container mt-5">Please log in to view your fines.</div>;
    if (loading) return <div className="container mt-5">Loading fines...</div>;

    return (
        <div className="container mt-4">
            <h1 className="mb-4">My Fines</h1>

            {error && (
                <div className="alert alert-danger">
                    {error}
                    <button className="btn-close float-end" onClick={() => setError(null)}></button>
                </div>
            )}

            {/* Create Fine Form */}
            <div className="mb-4">
                <h4>Create Fine</h4>
                <form onSubmit={createFine}>
                    <div className="row">
                        <div className="col-md-3">
                            <input className="form-control" name="bookId" placeholder="Book ID" required />
                        </div>
                        <div className="col-md-2">
                            <input className="form-control" name="daysLate" type="number" placeholder="Days Late" required />
                        </div>
                        <div className="col-md-2">
                            <input className="form-control" name="ratePerDay" type="number" step="0.01" placeholder="Rate Per Day" required />
                        </div>
                        <div className="col-md-2">
                            <button type="submit" className="btn btn-primary w-100">Create Fine</button>
                        </div>
                    </div>
                </form>
            </div>

            {/* Fines Table */}
            <div className="table-responsive">
                <table className="table table-bordered table-hover">
                    <thead className="table-dark">
                        <tr>
                            <th>Book ID</th>
                            <th>Days Late</th>
                            <th>Rate Per Day</th>
                            <th>Amount</th>
                            <th>Status</th>
                            <th>Created At</th>
                            <th>Actions</th>
                        </tr>
                    </thead>
                    <tbody>
                        {fines.length === 0 ? (
                            <tr>
                                <td colSpan="7" className="text-center">You don't have any fines.</td>
                            </tr>
                        ) : (
                            fines.map((fine) => (
                                <tr key={fine.fine_id}>
                                    <td>{fine.bookId}</td>
                                    <td>{fine.daysLate}</td>
                                    <td>₱{fine.ratePerDay.toFixed(2)}</td>
                                    <td>₱{fine.amount.toFixed(2)}</td>
                                    <td>
                                        <span className={`badge ${fine.amount > 0 ? "bg-danger" : "bg-success"}`}>
                                            {fine.amount > 0 ? "UNPAID" : "PAID"}
                                        </span>
                                    </td>
                                    <td>{new Date(fine.createdAt).toLocaleDateString()}</td>
                                    <td>
                                        <button
                                            className="btn btn-sm btn-outline-danger"
                                            onClick={() => deleteFine(fine.fine_id)}
                                        >
                                            Delete
                                        </button>
                                    </td>
                                </tr>
                            ))
                        )}
                    </tbody>
                </table>
            </div>
        </div>
    );
};

export default Fines;


