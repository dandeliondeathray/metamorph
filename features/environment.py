import pymetamorph.pymetamorph.metamorph as morph


def before_feature(context, feature):
    context.metamorph = morph.Metamorph()
