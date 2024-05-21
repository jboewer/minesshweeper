package game_test

import (
	"github.com/jboewer/minesshweeper/game"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func assertPlaceMine(t *testing.T, g *game.Game, x, y int) {
	t.Helper()
	assert.NoError(t, g.PlaceMine(x, y))
}

func TestItCanCreateGameWithSize(t *testing.T) {
	g, _ := game.New(10, 10)

	if g.GetGridWidth() != 10 {
		t.Errorf("Expected field width to be 10, got %d", g.GetGridWidth())
	}
	if g.GetGridHeight() != 10 {
		t.Errorf("Expected field width to be 10, got %d", g.GetGridHeight())
	}
}

func TestItCannotCreateAGameWithInvalidFieldSize(t *testing.T) {
	_, err1 := game.New(0, 0)
	_, err2 := game.New(1, 0)
	_, err3 := game.New(0, 1)
	_, err4 := game.New(-10, 1)
	_, err5 := game.New(1, -1)

	assert.ErrorIs(t, err1, game.ErrInvalidFieldSize)
	assert.ErrorIs(t, err2, game.ErrInvalidFieldSize)
	assert.ErrorIs(t, err3, game.ErrInvalidFieldSize)
	assert.ErrorIs(t, err4, game.ErrInvalidFieldSize)
	assert.ErrorIs(t, err5, game.ErrInvalidFieldSize)
}

func TestGame_PlaceMine(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		g, _ := game.New(10, 10)
		assert.NoError(t, g.PlaceMine(0, 0))
	})
	t.Run("OutOfBounds", func(t *testing.T) {
		g, _ := game.New(10, 10)

		assert.ErrorIs(t, g.PlaceMine(10, 0), game.ErrOutOfBounds)
		assert.ErrorIs(t, g.PlaceMine(0, 10), game.ErrOutOfBounds)
		assert.ErrorIs(t, g.PlaceMine(10, 10), game.ErrOutOfBounds)
		assert.ErrorIs(t, g.PlaceMine(-10, 0), game.ErrOutOfBounds)
		assert.ErrorIs(t, g.PlaceMine(0, -10), game.ErrOutOfBounds)
	})
	t.Run("Duplicate", func(t *testing.T) {
		g, _ := game.New(10, 10)
		assert.NoError(t, g.PlaceMine(0, 0))
		assert.ErrorIs(t, g.PlaceMine(0, 0), game.ErrDuplicateMine)
	})
}

func TestItGetsTheNumberOfMines(t *testing.T) {
	g, _ := game.New(10, 10)

	assertPlaceMine(t, g, 0, 0)
	assertPlaceMine(t, g, 1, 1)

	assert.Equal(t, 2, g.GetMineCount())
}

func TestItCanRevealCell(t *testing.T) {
	g, _ := game.New(10, 10)

	assert.Zero(t, g.RevealCell(0, 0))
}

func TestItCanRevealCellWithMine(t *testing.T) {
	g, _ := game.New(10, 10)

	assertPlaceMine(t, g, 0, 0)

	assert.Equal(t, -1, g.RevealCell(0, 0))
}

func TestGameStateIsPlayingWhenMinesLeftover(t *testing.T) {
	g, _ := game.New(10, 10)

	assertPlaceMine(t, g, 0, 0)

	assert.Equal(t, game.StatePlaying, g.State())
}

func TestGameStateIsLostWhenMineIsRevealed(t *testing.T) {
	g, _ := game.New(10, 10)

	assertPlaceMine(t, g, 0, 0)
	g.RevealCell(0, 0)

	assert.Equal(t, game.StateLost, g.State())
}

func TestItCanPlaceFlags(t *testing.T) {
	g, _ := game.New(10, 10)

	assert.NoError(t, g.PlaceFlag(0, 0))
}

func TestItCannotPlaceFlagOutOfBounds(t *testing.T) {
	g, _ := game.New(10, 10)

	assert.ErrorIs(t, g.PlaceFlag(10, 0), game.ErrOutOfBounds)
	assert.ErrorIs(t, g.PlaceFlag(0, 10), game.ErrOutOfBounds)
	assert.ErrorIs(t, g.PlaceFlag(10, 10), game.ErrOutOfBounds)
	assert.ErrorIs(t, g.PlaceFlag(-10, 0), game.ErrOutOfBounds)
	assert.ErrorIs(t, g.PlaceFlag(0, -10), game.ErrOutOfBounds)
}

func TestItGetsTheNumberOfPlacedFlags(t *testing.T) {
	g, _ := game.New(10, 10)

	assert.NoError(t, g.PlaceFlag(0, 0))
	assert.NoError(t, g.PlaceFlag(1, 1))

	assert.Equal(t, 2, g.GetFlagCount())
}

func TestGameIsWonWhenAllMinesCoveredWithFlag(t *testing.T) {
	g, _ := game.New(10, 10)

	assertPlaceMine(t, g, 0, 0)
	assertPlaceMine(t, g, 1, 1)

	assert.NoError(t, g.PlaceFlag(0, 0))
	assert.NoError(t, g.PlaceFlag(1, 1))

	assert.Equal(t, game.StateWon, g.State())
}

func TestGameIsNotWonIfTooManyFlags(t *testing.T) {
	g, _ := game.New(10, 10)

	assertPlaceMine(t, g, 0, 0)
	assertPlaceMine(t, g, 1, 1)

	assert.NoError(t, g.PlaceFlag(0, 0))
	assert.NoError(t, g.PlaceFlag(1, 1))
	assert.NoError(t, g.PlaceFlag(2, 2))

	assert.Equal(t, game.StatePlaying, g.State())
}

func TestGame_PlaceMines_ValidSizes(t *testing.T) {
	t.Run("1x1", func(t *testing.T) {
		g, _ := game.New(1, 1)
		m := game.Grid{
			{0},
		}
		assert.NoError(t, g.PlaceMines(m))
	})
	t.Run("5x1", func(t *testing.T) {
		g, _ := game.New(5, 1)
		m := game.Grid{
			{0, 0, 0, 0, 0},
		}
		assert.NoError(t, g.PlaceMines(m))
	})
	t.Run("1x5", func(t *testing.T) {
		g, _ := game.New(1, 5)
		m := game.Grid{
			{0},
			{0},
			{0},
			{0},
			{0},
		}
		assert.NoError(t, g.PlaceMines(m))
	})
}

func TestGame_PlaceMines_InvalidSizes(t *testing.T) {
	t.Run("mismatch width", func(t *testing.T) {
		g, _ := game.New(3, 5)
		m := game.Grid{
			{0},
			{0},
			{0},
			{0},
			{0},
		}
		assert.ErrorIs(t, g.PlaceMines(m), game.ErrInvalidFieldSize)
	})
	t.Run("mismatch height", func(t *testing.T) {
		g, _ := game.New(5, 3)
		m := game.Grid{
			{0, 0, 0, 0, 0},
		}
		assert.ErrorIs(t, g.PlaceMines(m), game.ErrInvalidFieldSize)
	})
}

func TestGame_PlaceMines(t *testing.T) {
	t.Run("Diagonal", func(t *testing.T) {
		g, _ := game.New(5, 5)
		m := game.Grid{
			{1, 0, 0, 0, 0},
			{0, 1, 0, 0, 0},
			{0, 0, 1, 0, 0},
			{0, 0, 0, 1, 0},
			{0, 0, 0, 0, 1},
		}
		assert.NoError(t, g.PlaceMines(m))
		assert.Equal(t, 5, g.GetMineCount())
		assert.Equal(t, -1, g.RevealCell(0, 0))
		assert.Equal(t, -1, g.RevealCell(1, 1))
		assert.Equal(t, -1, g.RevealCell(2, 2))
		assert.Equal(t, -1, g.RevealCell(3, 3))
		assert.Equal(t, -1, g.RevealCell(4, 4))
	})
	t.Run("Bottom/left", func(t *testing.T) {
		g, _ := game.New(5, 5)
		m := game.Grid{
			{0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0},
			{1, 0, 0, 0, 0},
		}
		assert.NoError(t, g.PlaceMines(m))
		assert.Equal(t, 1, g.GetMineCount())
		assert.Equal(t, -1, g.RevealCell(0, 4))
	})
}

func TestGame_NewFromPlacementMap(t *testing.T) {
	m := game.Grid{
		{1, 0, 0, 0, 0},
		{0, 1, 0, 0, 0},
		{0, 0, 1, 0, 0},
		{0, 0, 0, 1, 0},
		{0, 0, 0, 0, 1},
	}

	g, err := game.NewFromGrid(m)
	assert.NoError(t, err)

	assert.Equal(t, 5, g.GetMineCount())
	assert.Equal(t, -1, g.RevealCell(0, 0))
	assert.Equal(t, -1, g.RevealCell(1, 1))
	assert.Equal(t, -1, g.RevealCell(2, 2))
	assert.Equal(t, -1, g.RevealCell(3, 3))
	assert.Equal(t, -1, g.RevealCell(4, 4))
}

func TestMinePlacementMap_Validate(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		m := game.Grid{
			{1, 0, 0, 0, 0},
			{0, 1, 0, 0, 0},
		}
		assert.NoError(t, m.Validate())
	})
	t.Run("Invalid", func(t *testing.T) {
		m := game.Grid{
			{1, 0, 0, 0, 0},
			{0, 1, 0, 0},
		}
		assert.ErrorIs(t, m.Validate(), game.ErrMalformedGrid)
	})
}

func TestGame_RevealCell_ItReturnsTheNumberOfAdjacentMines(t *testing.T) {
	t.Run("No mines", func(t *testing.T) {
		g, _ := game.New(5, 5)

		assert.Equal(t, 0, g.RevealCell(0, 0))
	})

	t.Run("One mine", func(t *testing.T) {
		g, _ := game.NewFromGrid(game.Grid{
			{1, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
		})

		assert.Equal(t, 1, g.RevealCell(1, 1))
	})

	t.Run("All mines but the revealed", func(t *testing.T) {
		g, _ := game.NewFromGrid(game.Grid{
			{1, 1, 1},
			{1, 0, 1},
			{1, 1, 1},
		})

		assert.Equal(t, 8, g.RevealCell(1, 1))
	})
	t.Run("top/left", func(t *testing.T) {
		g, _ := game.NewFromGrid(game.Grid{
			{0, 1, 1},
			{1, 1, 1},
			{1, 1, 1},
		})

		assert.Equal(t, 3, g.RevealCell(0, 0))
	})
	t.Run("top/right", func(t *testing.T) {
		g, _ := game.NewFromGrid(game.Grid{
			{1, 1, 0},
			{1, 1, 1},
			{1, 1, 1},
		})

		assert.Equal(t, 3, g.RevealCell(2, 0))
	})
	t.Run("bottom/left", func(t *testing.T) {
		g, _ := game.NewFromGrid(game.Grid{
			{1, 1, 1},
			{1, 1, 1},
			{0, 1, 1},
		})

		assert.Equal(t, 3, g.RevealCell(0, 2))
	})

	t.Run("flat left", func(t *testing.T) {
		g, _ := game.NewFromGrid(game.Grid{
			{0, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		})

		assert.Equal(t, 1, g.RevealCell(0, 0))
	})
}

func TestGame_GetGrid(t *testing.T) {
	t.Run("No cells revealed", func(t *testing.T) {
		g, _ := game.New(5, 1)

		expected := game.Grid{
			{game.CellUnrevealed, game.CellUnrevealed, game.CellUnrevealed, game.CellUnrevealed, game.CellUnrevealed},
		}
		assertEqualGrid(t, expected, g.GetGrid())
	})
	t.Run("All cells revealed", func(t *testing.T) {
		g, _ := game.New(5, 1)
		g.RevealCell(0, 0)
		g.RevealCell(1, 0)
		g.RevealCell(2, 0)
		g.RevealCell(3, 0)
		g.RevealCell(4, 0)

		expected := game.Grid{
			{0, 0, 0, 0, 0},
		}
		assertEqualGrid(t, expected, g.GetGrid())
	})
	t.Run("flagged cell", func(t *testing.T) {
		g, _ := game.New(5, 1)
		g.PlaceFlag(1, 0)

		expected := game.Grid{
			{game.CellUnrevealed, game.CellFlag, game.CellUnrevealed, game.CellUnrevealed, game.CellUnrevealed},
		}
		assertEqualGrid(t, expected, g.GetGrid())
	})
}

func TestGame_GetGrid_WithMines(t *testing.T) {
	g, _ := game.New(5, 1)
	g.PlaceMine(0, 0)

	expected := game.Grid{
		{game.CellMine, game.CellUnrevealed, game.CellUnrevealed, game.CellUnrevealed, game.CellUnrevealed},
	}

	g.RevealCell(0, 0)

	assertEqualGrid(t, expected, g.GetGrid())
}

func TestAFlagCanBeRemoved(t *testing.T) {
	g, _ := game.New(5, 1)
	g.PlaceFlag(1, 0)
	g.RemoveFlag(1, 0)

	expected := game.Grid{
		{game.CellUnrevealed, game.CellUnrevealed, game.CellUnrevealed, game.CellUnrevealed, game.CellUnrevealed},
	}
	assertEqualGrid(t, expected, g.GetGrid())
}

func TestItRevealsAdjacentCellsIfCellHasNone(t *testing.T) {
	g, _ := game.NewFromGrid(game.Grid{
		{1, 1, 1, 1, 1},
		{1, 0, 0, 0, 1},
		{1, 0, 0, 0, 1},
		{1, 0, 0, 0, 1},
		{1, 1, 1, 1, 1},
	})

	g.RevealCell(2, 2)

	expected := game.Grid{
		{game.CellUnrevealed, game.CellUnrevealed, game.CellUnrevealed, game.CellUnrevealed, game.CellUnrevealed},
		{game.CellUnrevealed, 5, 3, 5, game.CellUnrevealed},
		{game.CellUnrevealed, 3, 0, 3, game.CellUnrevealed},
		{game.CellUnrevealed, 5, 3, 5, game.CellUnrevealed},
		{game.CellUnrevealed, game.CellUnrevealed, game.CellUnrevealed, game.CellUnrevealed, game.CellUnrevealed},
	}
	assertEqualGrid(t, expected, g.GetGrid())
}

func assertEqualGrid(t *testing.T, expected, actual game.Grid) {
	t.Helper()

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}
