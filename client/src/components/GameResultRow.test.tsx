import { render, screen } from "@testing-library/react"
import userEvent from "@testing-library/user-event"
import GameResultRow from "./GameResultRow"
import type { Games } from "../types/types"
import { vi } from "vitest"

const mockNavigate = vi.fn()

vi.mock("@tanstack/react-router", () => ({
  useNavigate: () => mockNavigate,
}))

describe("<GameResultRow />", () => {
  const game: Games = {
    id: 1,
    guid: "abc123",
    name: "Mock Game",
    image: {
      medium_url: "",
      icon_url: "",
      tiny_url: "",
      thumb_url: "https://example.com/thumb.jpg",
      small_url: "",
      super_url: "",
    },
    site_detail_url: "https://example.com/game",
  }

  test("clicking the row triggers navigation with correct route", async () => {
    render(<GameResultRow game={game} />)
    const user = userEvent.setup()

    const row = await screen.findByText("Mock Game")
    await user.click(row)

    expect(mockNavigate).toHaveBeenCalledWith({
      to: "/games/$guid",
      params: { guid: "abc123" },
    })
  })
})
