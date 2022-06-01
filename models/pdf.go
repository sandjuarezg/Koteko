package models

import (
	"strconv"

	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
)

func GenerateBuyPDF(product Producto, cant int) (err error) {
	m := pdf.NewMaroto(consts.Portrait, consts.Letter)

	m.Row(40, func() {
		m.Col(12, func() {
			_ = m.FileImage("./public/img/koteko/logo.png", props.Rect{
				Center:  true,
				Percent: 60,
			})
		})
	})

	m.Row(20, func() {
		m.Col(12, func() {
			m.Text("Gracias por tu compra, con tu preferencia nos ayudas a cumplir nuestras metas.", props.Text{
				Align: consts.Center,
				Size:  15,
				Top:   14,
			})
		})
	})

	m.Line(10)

	m.Row(20, func() {
		m.Col(12, func() {
			m.Text(product.Producto, props.Text{
				Size: 15,
			})
		})
	})

	m.Row(50, func() {
		m.Col(12, func() {
			m.Text(product.Descripcion, props.Text{
				Size: 15,
			})
		})
	})

	m.SetBorder(true)

	m.Row(10, func() {
		m.Col(4, func() {
			m.Text("Precio", props.Text{
				Style: consts.Bold,
				Size:  18,
				Align: consts.Center,
			})
		})
		m.Col(4, func() {
			m.Text("Cantidad", props.Text{
				Style: consts.Bold,
				Size:  18,
				Align: consts.Center,
			})
		})
		m.Col(4, func() {
			m.Text("Total", props.Text{
				Style: consts.Bold,
				Size:  18,
				Align: consts.Center,
			})
		})
	})

	m.Row(10, func() {
		m.Col(4, func() {
			m.Text(strconv.Itoa(int(product.Precio)), props.Text{
				Top:   2,
				Size:  14,
				Align: consts.Center,
			})
		})
		m.Col(4, func() {
			m.Text(strconv.Itoa(int(cant)), props.Text{
				Top:   2,
				Size:  14,
				Align: consts.Center,
			})
		})
		m.Col(4, func() {
			m.Text(strconv.Itoa(int(int(product.Precio)*cant)), props.Text{
				Top:   2,
				Size:  14,
				Align: consts.Center,
			})
		})
	})

	m.SetBorder(false)

	err = m.OutputFileAndClose("/home/sand/Desktop/comprobantePago.pdf")
	if err != nil {
		return
	}

	return
}
